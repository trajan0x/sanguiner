import { useState, useEffect } from 'react'
import { useAccount } from 'wagmi'
import { Address, multicall, erc20ABI, getAccount } from '@wagmi/core'
import { BRIDGABLE_TOKENS } from '@/constants/tokens'
import { Token } from '../types'
import { AddressZero } from '@ethersproject/constants'
import multicallABI from '@/constants/abis/multicall.json'
import { getSortedBridgableTokens } from '../actions/getSortedBridgableTokens'
import { ChainId } from '@/constants/chains'
import { sortByTokenBalance } from '../sortTokens'
import { fetchBalance } from '@wagmi/core'

//move to constants file later
const MULTICALL3_ADDRESS: Address = '0xcA11bde05977b3631167028862bE2a173976CA11'

export const getTokensByChainId = async (chainId: number) => {
  const { address } = getAccount()

  const tokens = BRIDGABLE_TOKENS[chainId]

  return await sortByTokenBalance(tokens, chainId, address)
}

export const usePortfolioBalances = () => {
  const { address } = getAccount()
  const availableChains = Object.keys(BRIDGABLE_TOKENS)

  useEffect(() => {
    const getData = async () => {
      availableChains.forEach(async (chainId) => {
        const response = await getTokensByChainId(Number(chainId))
        console.log('response chainId: ', chainId, response)
      })
    }
    getData()
  }, [])
}

const useTokenApprovals = () => {}
