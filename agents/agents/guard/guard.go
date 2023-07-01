package guard

import (
	"context"
	"fmt"
	"github.com/synapsecns/sanguine/core/metrics"
	signerConfig "github.com/synapsecns/sanguine/ethergo/signer/config"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"

	"github.com/synapsecns/sanguine/agents/config"
	"github.com/synapsecns/sanguine/agents/domains"
	"github.com/synapsecns/sanguine/agents/domains/evm"
	"github.com/synapsecns/sanguine/agents/types"
	"github.com/synapsecns/sanguine/ethergo/signer/signer"
)

// Guard scans origins for latest state and submits snapshots to the Summit.
type Guard struct {
	bondedSigner       signer.Signer
	unbondedSigner     signer.Signer
	domains            []domains.DomainClient
	summitDomain       domains.DomainClient
	refreshInterval    time.Duration
	summitLatestStates map[uint32]types.State
	// TODO: change to metrics type
	originLatestStates map[uint32]types.State
	handler            metrics.Handler
}

// NewGuard creates a new guard.
//
//nolint:cyclop
func NewGuard(ctx context.Context, cfg config.AgentConfig, handler metrics.Handler) (_ Guard, err error) {
	guard := Guard{
		refreshInterval: time.Second * time.Duration(cfg.RefreshIntervalSeconds),
	}
	guard.domains = []domains.DomainClient{}

	guard.bondedSigner, err = signerConfig.SignerFromConfig(ctx, cfg.BondedSigner)
	if err != nil {
		return Guard{}, fmt.Errorf("error with bondedSigner, could not create guard: %w", err)
	}

	guard.unbondedSigner, err = signerConfig.SignerFromConfig(ctx, cfg.UnbondedSigner)
	if err != nil {
		return Guard{}, fmt.Errorf("error with unbondedSigner, could not create guard: %w", err)
	}

	for domainName, domain := range cfg.Domains {
		var domainClient domains.DomainClient

		chainRPCURL := fmt.Sprintf("%s/confirmations/1/rpc/%d", cfg.BaseOmnirpcURL, domain.DomainID)
		domainClient, err = evm.NewEVM(ctx, domainName, domain, chainRPCURL)
		if err != nil {
			return Guard{}, fmt.Errorf("failing to create evm for domain, could not create guard for: %w", err)
		}
		guard.domains = append(guard.domains, domainClient)
		if domain.DomainID == cfg.SummitDomainID {
			guard.summitDomain = domainClient
		}
	}

	guard.summitLatestStates = make(map[uint32]types.State, len(guard.domains))
	guard.originLatestStates = make(map[uint32]types.State, len(guard.domains))

	guard.handler = handler

	return guard, nil
}

//nolint:cyclop
func (g Guard) loadSummitLatestStates(parentCtx context.Context) {
	for _, domain := range g.domains {
		ctx, span := g.handler.Tracer().Start(parentCtx, "loadSummitLatestStates", trace.WithAttributes(
			attribute.Int("domain", int(domain.Config().DomainID)),
		))

		originID := domain.Config().DomainID
		latestState, err := g.summitDomain.Summit().GetLatestAgentState(ctx, originID, g.bondedSigner)
		if err != nil {
			latestState = nil
			logger.Errorf("Failed calling GetLatestAgentState for originID %d on the Summit contract: err = %v", originID, err)
			span.AddEvent("Failed calling GetLatestAgentState for originID on the Summit contract", trace.WithAttributes(
				attribute.Int("originID", int(originID)),
				attribute.String("err", err.Error()),
			))
		}
		if latestState != nil && latestState.Nonce() > uint32(0) {
			g.summitLatestStates[originID] = latestState
		}

		span.End()
	}
}

//nolint:cyclop
func (g Guard) loadOriginLatestStates(parentCtx context.Context) {
	for _, domain := range g.domains {
		ctx, span := g.handler.Tracer().Start(parentCtx, "loadOriginLatestStates", trace.WithAttributes(
			attribute.Int("domain", int(domain.Config().DomainID)),
		))

		originID := domain.Config().DomainID
		latestState, err := domain.Origin().SuggestLatestState(ctx)
		if err != nil {
			latestState = nil
			logger.Errorf("Failed calling SuggestLatestState for originID %d on the Origin contract: %v", originID, err)
			span.AddEvent("Failed calling SuggestLatestState for originID on the Origin contract", trace.WithAttributes(
				attribute.Int("originID", int(originID)),
				attribute.String("err", err.Error()),
			))
		} else if latestState == nil || latestState.Nonce() == uint32(0) {
			logger.Errorf("No latest state found for origin id %d", originID)
			span.AddEvent("No latest state found for origin id", trace.WithAttributes(
				attribute.Int("originID", int(originID)),
			))
		}
		if latestState != nil {
			// TODO: if overwriting, end span and start a new one
			g.originLatestStates[originID] = latestState
		}

		span.End()
	}
}

//nolint:cyclop
func (g Guard) getLatestSnapshot() (types.Snapshot, map[uint32]types.State) {
	statesToSubmit := make(map[uint32]types.State, len(g.domains))
	for _, domain := range g.domains {
		originID := domain.Config().DomainID
		summitLatest, ok := g.summitLatestStates[originID]
		if !ok || summitLatest == nil || summitLatest.Nonce() == 0 {
			summitLatest = nil
		}
		originLatest, ok := g.originLatestStates[originID]
		if !ok || originLatest == nil || originLatest.Nonce() == 0 {
			continue
		}
		if summitLatest != nil && summitLatest.Nonce() >= originLatest.Nonce() {
			// Here this guard already submitted this state
			continue
		}
		// TODO: add event for submitting that state
		statesToSubmit[originID] = originLatest
	}
	snapshotStates := make([]types.State, 0, len(statesToSubmit))
	for _, state := range statesToSubmit {
		if state.Nonce() == 0 {
			continue
		}
		snapshotStates = append(snapshotStates, state)
	}
	if len(snapshotStates) > 0 {
		return types.NewSnapshot(snapshotStates), statesToSubmit
	}
	//nolint:nilnil
	return nil, nil
}

//nolint:cyclop
func (g Guard) submitLatestSnapshot(parentCtx context.Context) {
	ctx, span := g.handler.Tracer().Start(parentCtx, "submitLatestSnapshot", trace.WithAttributes(
		attribute.Int("domain", int(g.summitDomain.Config().DomainID)),
	))

	defer func() {
		span.End()
	}()

	snapshot, statesToSubmit := g.getLatestSnapshot()
	if snapshot == nil {
		return
	}

	snapshotSignature, encodedSnapshot, _, err := snapshot.SignSnapshot(ctx, g.bondedSigner)
	if err != nil {
		logger.Errorf("Error signing snapshot: %v", err)
		span.AddEvent("Error signing snapshot", trace.WithAttributes(
			attribute.String("err", err.Error()),
		))
	} else {
		err = g.summitDomain.Inbox().SubmitSnapshot(ctx, g.unbondedSigner, encodedSnapshot, snapshotSignature)
		if err != nil {
			logger.Errorf("Failed to submit snapshot to inbox: %v", err)
			span.AddEvent("Failed to submit snapshot to inbox", trace.WithAttributes(
				attribute.String("err", err.Error()),
			))
		} else {
			for originID, state := range statesToSubmit {
				g.summitLatestStates[originID] = state
			}
		}
	}
}

//nolint:cyclop
func (g Guard) loadNumAttestations(parentCtx context.Context) {
	for _, domain := range g.domains {
		ctx, span := g.handler.Tracer().Start(parentCtx, "loadNumAttestations", trace.WithAttributes(
			attribute.Int("domain", int(domain.Config().DomainID)),
		))

		if domain.Config().DomainID == g.summitDomain.Config().DomainID {
			continue
		}

		numAttestations, err := domain.Destination().AttestationsAmount(ctx)
		if err != nil {
			logger.Errorf("Failed calling AttestationsAmount for destinationID %d on the Destination contract: err = %v", domain.Config().DomainID, err)
			span.AddEvent("Failed calling AttestationsAmount for destinationID on the Destination contract", trace.WithAttributes(
				attribute.Int("destinationID", int(domain.Config().DomainID)),
				attribute.String("err", err.Error()),
			))
		}

		if numAttestations > 0 {
			fmt.Printf("CRONIN destinationID %d domain.Config().DomainID has an Attestation count of %v\n", domain.Config().DomainID, numAttestations)
			for count := numAttestations; count > 0; count-- {
				i := count - 1
				if count == 3 {
					fmt.Printf("CRONIN count == 3")
				}
				attestation, attestationSignature, err := domain.Destination().GetAttestation(ctx, i)
				if err != nil {
					logger.Errorf("Failed calling GetAttestation for destinationID %d and index %v on the Destination contract: err = %v", domain.Config().DomainID, i, err)
					span.AddEvent("Failed calling GetAttestation for destinationID and index on the Destination contract", trace.WithAttributes(
						attribute.Int("destinationID", int(domain.Config().DomainID)),
						attribute.Int64("index", int64(i)),
						attribute.String("err", err.Error()),
					))
					continue
				}
				if attestation == nil || attestation.Nonce() == 0 {
					continue
				}

				isValidAttestation, err := g.summitDomain.Summit().IsValidAttestation(ctx, attestation)
				if err != nil {
					// TODO (joe): This could be very serious if we can't reach Summit for an extended period of time.
					// We need to monitor and alert cases when an attestation can't be checked.
					// If need be, we should manually pause the destination from using this attestation.
					logger.Errorf("Failed calling IsValidAttestation for destinationID %d and index %d on the Summit contract: err = %v", domain.Config().DomainID, i, err)
					span.AddEvent("Failed calling IsValidAttestation for destinationID and index on the Summit contract", trace.WithAttributes(
						attribute.Int("destinationID", int(domain.Config().DomainID)),
						attribute.Int64("index", int64(i)),
						attribute.String("err", err.Error()),
					))
					continue
				}
				if isValidAttestation {
					fmt.Printf("CRONIN Attestation is valid!!!\n")
				} else {
					// TODO (joe): Submit fraud report here
					fmt.Printf("CRONIN Attestation is NOT valid!!! WE FOUND FRAUD!!!!!!\n")

					// First thing is to call verify on Summit
					attPayload, err := types.EncodeAttestation(attestation)
					if err != nil {
						logger.Errorf("Failed EncodeAttestation for destinationID %d and index %d: err = %v", domain.Config().DomainID, i, err)
						span.AddEvent("Failed calling EncodeAttestation for destinationID and index", trace.WithAttributes(
							attribute.Int("destinationID", int(domain.Config().DomainID)),
							attribute.Int64("index", int64(i)),
							attribute.String("err", err.Error()),
						))
						continue
					}
					attSignature, err := types.EncodeSignature(attestationSignature)
					if err != nil {
						logger.Errorf("Failed EncodeSignature for destinationID %d and index %d: err = %v", domain.Config().DomainID, i, err)
						span.AddEvent("Failed calling EncodeSignature for destinationID and index", trace.WithAttributes(
							attribute.Int("destinationID", int(domain.Config().DomainID)),
							attribute.Int64("index", int64(i)),
							attribute.String("err", err.Error()),
						))
						continue
					}

					err = g.summitDomain.Inbox().VerifyAttestation(ctx, g.unbondedSigner, attPayload, attSignature)
					if err != nil {
						logger.Errorf("Failed to call VerifyAttestation on Inbox: %v", err)
						span.AddEvent("Failed to call VerifyAttestation on inbox", trace.WithAttributes(
							attribute.String("err", err.Error()),
						))
						continue
					}

					// Then submit fraud report on destination
				}
			}
		}

		span.End()
	}
}

// Start starts the guard.
//
//nolint:cyclop
func (g Guard) Start(ctx context.Context) error {
	// First initialize a map to track what was the last state signed by this guard
	g.loadSummitLatestStates(ctx)

	for {
		select {
		// parent loop terminated
		case <-ctx.Done():
			logger.Info("Guard exiting without error")
			return nil
		case <-time.After(g.refreshInterval):
			g.loadOriginLatestStates(ctx)
			g.submitLatestSnapshot(ctx)
			g.loadNumAttestations(ctx)
		}
	}
}
