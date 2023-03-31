import '@styles/global.css'
import '@rainbow-me/rainbowkit/styles.css'
import type { AppProps } from 'next/app'
import { Provider as EthersProvider } from '@ethersproject/abstract-provider'
import { JsonRpcProvider } from '@ethersproject/providers'
import {
  boba,
  cronos,
  dfk,
  dogechain,
  klaytn,
} from '@constants/extraWagmiChains'
import { WagmiConfig, configureChains, createClient } from 'wagmi'
import {
  arbitrum,
  aurora,
  avalanche,
  bsc,
  canto,
  celo,
  fantom,
  harmonyOne,
  mainnet,
  metis,
  moonbeam,
  moonriver,
  optimism,
  polygon,
} from 'wagmi/chains'
import {
  RainbowKitProvider,
  darkTheme,
  getDefaultWallets,
} from '@rainbow-me/rainbowkit'
import { alchemyProvider } from 'wagmi/providers/alchemy'
import { publicProvider } from 'wagmi/providers/public'
import { CHAIN_INFO_MAP } from '@constants/networks'

import { SynapseProvider } from '@/utils/SynapseProvider'
export default function App({ Component, pageProps }: AppProps) {
  const rawChains = [
    mainnet,
    arbitrum,
    aurora,
    avalanche,
    bsc,
    canto,
    fantom,
    harmonyOne,
    metis,
    moonbeam,
    moonriver,
    optimism,
    polygon,
    klaytn,
    cronos,
    dfk,
    dogechain,
    boba,
  ]

  // Add custom icons
  const chainsWithIcons: any[] = []
  for (const chain of rawChains) {
    chainsWithIcons.push({
      ...chain,
      iconUrl: CHAIN_INFO_MAP[chain.id].chainImg.src,
    })
  }
  const { chains, provider } = configureChains(chainsWithIcons, [
    alchemyProvider({ apiKey: '_UFN4P3jhI9zYma6APzoKX5aqKKadp2V' }),
    publicProvider(),
  ])

  const { connectors } = getDefaultWallets({
    appName: 'Synapse',
    chains,
  })

  const wagmiClient = createClient({
    autoConnect: true,
    connectors,
    provider,
  })

  // Synapse client
  const synapseProviders: EthersProvider[] = []
  chains.map((chain) => {
    const rpc: EthersProvider = new JsonRpcProvider(
      chain.rpcUrls.default.http[0]
    )
    synapseProviders.push(rpc)
  })
  return (
    <WagmiConfig client={wagmiClient}>
      <RainbowKitProvider chains={chains} theme={darkTheme()}>
        <SynapseProvider
          chainIds={chains.map((chain) => chain.id)}
          providers={synapseProviders}
        >
          <Component {...pageProps} />
        </SynapseProvider>
      </RainbowKitProvider>
    </WagmiConfig>
  )
}
