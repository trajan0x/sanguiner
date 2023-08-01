import _ from 'lodash'

import { EXISTING_BRIDGE_ROUTES } from '@/constants/existing-bridge-routes'
import { RouteQueryFields } from './generateRoutePossibilities'
import { getTokenAndChainId } from './getTokenAndChainId'

export const getToChainIds = ({
  fromChainId,
  fromTokenRouteSymbol,
  toChainId,
  toTokenRouteSymbol,
}: RouteQueryFields) => {
  if (
    fromChainId === null &&
    fromTokenRouteSymbol === null &&
    toChainId === null &&
    toTokenRouteSymbol === null
  ) {
    return _(EXISTING_BRIDGE_ROUTES)
      .values()
      .flatten()
      .map((token) => getTokenAndChainId(token).chainId)
      .uniq()
      .value()
  }

  if (
    fromChainId &&
    fromTokenRouteSymbol === null &&
    toChainId === null &&
    toTokenRouteSymbol === null
  ) {
    return _(EXISTING_BRIDGE_ROUTES)
      .entries()
      .filter(([key, _values]) => {
        const { chainId } = getTokenAndChainId(key)
        return chainId === fromChainId
      })
      .map(([_key, values]) => values)
      .flatten()
      .map((token) => getTokenAndChainId(token).chainId)
      .uniq()
      .value()
  }

  if (
    fromChainId === null &&
    fromTokenRouteSymbol &&
    toChainId === null &&
    toTokenRouteSymbol === null
  ) {
    return _(EXISTING_BRIDGE_ROUTES)
      .entries()
      .filter(([key, _values]) => {
        const { symbol } = getTokenAndChainId(key)
        return symbol === fromTokenRouteSymbol
      })
      .map(([_key, values]) => values)
      .flatten()
      .map((token) => getTokenAndChainId(token).chainId)
      .uniq()
      .value()
  }

  if (
    fromChainId &&
    fromTokenRouteSymbol &&
    toChainId === null &&
    toTokenRouteSymbol === null
  ) {
    return _.uniq(
      EXISTING_BRIDGE_ROUTES[`${fromTokenRouteSymbol}-${fromChainId}`].map(
        (token) => getTokenAndChainId(token).chainId
      )
    )
  }

  if (
    fromChainId === null &&
    fromTokenRouteSymbol === null &&
    toChainId &&
    toTokenRouteSymbol === null
  ) {
    return _(EXISTING_BRIDGE_ROUTES)
      .values()
      .flatten()
      .map((token) => getTokenAndChainId(token).chainId)
      .uniq()
      .value()
  }

  if (
    fromChainId &&
    fromTokenRouteSymbol === null &&
    toChainId &&
    toTokenRouteSymbol === null
  ) {
    return _(EXISTING_BRIDGE_ROUTES)
      .entries()
      .filter(([key, _values]) => key.endsWith(`-${fromChainId}`))
      .map(([_key, values]) => values)
      .flatten()
      .map((token) => getTokenAndChainId(token).chainId)
      .uniq()
      .value()
  }

  if (
    fromChainId === null &&
    fromTokenRouteSymbol &&
    toChainId &&
    toTokenRouteSymbol === null
  ) {
    return _(EXISTING_BRIDGE_ROUTES)
      .entries()
      .filter(([key, _values]) => key.startsWith(`${fromTokenRouteSymbol}-`))
      .filter(([_key, values]) =>
        values.some((v) => getTokenAndChainId(v).chainId === toChainId)
      )
      .map(([_key, values]) => values)
      .flatten()
      .filter((token) => token.endsWith(`-${toChainId}`))
      .map((token) => getTokenAndChainId(token).chainId)
      .uniq()
      .value()
  }

  if (
    fromChainId &&
    fromTokenRouteSymbol &&
    toChainId &&
    toTokenRouteSymbol === null
  ) {
    return _.uniq(
      EXISTING_BRIDGE_ROUTES[`${fromTokenRouteSymbol}-${fromChainId}`].map(
        (token) => getTokenAndChainId(token).chainId
      )
    )
  }

  if (
    fromChainId === null &&
    fromTokenRouteSymbol === null &&
    toChainId === null &&
    toTokenRouteSymbol
  ) {
    return _(EXISTING_BRIDGE_ROUTES)
      .values()
      .flatten()
      .filter((token) => token.startsWith(`${toTokenRouteSymbol}-`))
      .map((token) => getTokenAndChainId(token).chainId)
      .uniq()
      .value()
  }

  if (
    fromChainId &&
    fromTokenRouteSymbol === null &&
    toChainId === null &&
    toTokenRouteSymbol
  ) {
    return _(EXISTING_BRIDGE_ROUTES)
      .entries()
      .filter(([key, _values]) => key.endsWith(`-${fromChainId}`))
      .map(([_key, values]) => values)
      .flatten()
      .filter((token) => token.startsWith(`${toTokenRouteSymbol}-`))
      .map((token) => getTokenAndChainId(token).chainId)
      .uniq()
      .value()
  }

  if (
    fromChainId === null &&
    fromTokenRouteSymbol &&
    toChainId === null &&
    toTokenRouteSymbol
  ) {
    return _(EXISTING_BRIDGE_ROUTES)
      .entries()
      .filter(([key, _values]) => key.startsWith(`${fromTokenRouteSymbol}-`))
      .map(([_key, values]) => values)
      .flatten()
      .filter((token) => token.startsWith(`${toTokenRouteSymbol}-`))
      .map((token) => getTokenAndChainId(token).chainId)
      .uniq()
      .value()
  }

  if (
    fromChainId &&
    fromTokenRouteSymbol &&
    toChainId === null &&
    toTokenRouteSymbol
  ) {
    return _.uniq(
      EXISTING_BRIDGE_ROUTES[`${fromTokenRouteSymbol}-${fromChainId}`].map(
        (token) => getTokenAndChainId(token).chainId
      )
    )
  }

  if (
    fromChainId === null &&
    fromTokenRouteSymbol === null &&
    toChainId &&
    toTokenRouteSymbol
  ) {
    return _(EXISTING_BRIDGE_ROUTES)
      .values()
      .flatten()
      .filter((token) => token.startsWith(`${toTokenRouteSymbol}`))
      .map((token) => getTokenAndChainId(token).chainId)
      .uniq()
      .value()
  }

  if (
    fromChainId &&
    fromTokenRouteSymbol === null &&
    toChainId &&
    toTokenRouteSymbol
  ) {
    return _(EXISTING_BRIDGE_ROUTES)
      .entries()
      .filter(([key, _values]) => key.endsWith(`-${fromChainId}`))
      .map(([_key, values]) => values)
      .flatten()
      .filter((token) => token === `${toTokenRouteSymbol}-${toChainId}`)
      .map((token) => getTokenAndChainId(token).chainId)
      .uniq()
      .value()
  }

  if (
    fromChainId === null &&
    fromTokenRouteSymbol &&
    toChainId &&
    toTokenRouteSymbol
  ) {
    return _(EXISTING_BRIDGE_ROUTES)
      .entries()
      .filter(([key, _values]) => key.startsWith(`${fromTokenRouteSymbol}-`))
      .map(([_key, values]) => values)
      .flatten()
      .filter((token) => token === `${toTokenRouteSymbol}-${toChainId}`)
      .map((token) => getTokenAndChainId(token).chainId)
      .uniq()
      .value()
  }

  if (fromChainId && fromTokenRouteSymbol && toChainId && toTokenRouteSymbol) {
    return _.uniq(
      EXISTING_BRIDGE_ROUTES[`${fromTokenRouteSymbol}-${fromChainId}`].map(
        (token) => getTokenAndChainId(token).chainId
      )
    )
  }
}
