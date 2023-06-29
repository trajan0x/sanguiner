import { useState, useEffect } from 'react'
import { useAccount } from 'wagmi'
import { Address, multicall, erc20ABI, getAccount } from '@wagmi/core'
import { BRIDGABLE_TOKENS } from '@/constants/tokens'
import { Token } from '../types'
import { sortByTokenBalance, TokenAndBalance } from '../sortTokens'
import { BigNumber } from 'ethers'
import { AddressZero } from '@ethersproject/constants'

const ROUTER_ADDRESS = '0x7E7A0e201FD38d3ADAA9523Da6C109a07118C96a'
interface NetworkTokenBalances {
  [index: number]: TokenAndBalance[]
}

export const getTokensByChainId = async (
  owner: string,
  tokens: Token[],
  chainId: number
): Promise<TokenAndBalance[]> => {
  return await sortByTokenBalance(tokens, chainId, owner)
}

interface TokenWithBalanceAndAllowance {
  token: Token
  balance: BigNumber
  allowance: BigNumber
}

function mergeBalancesAndAllowances(
  balances: { token: Token; balance: BigNumber }[],
  allowances: { token: Token; allowance: BigNumber }[]
): TokenWithBalanceAndAllowance[] {
  return balances.map((balance) => {
    const correspondingAllowance = allowances.find(
      (item2) => item2.token === balance.token
    )

    if (correspondingAllowance) {
      return {
        token: balance.token,
        balance: balance.balance,
        allowance: correspondingAllowance.allowance,
      }
    }

    // if no allowance is matched with corresponding balance
    // e.g native gas tokens
    return {
      token: balance.token,
      balance: balance.balance,
      allowance: null,
    }
  })
}

export const usePortfolioBalances = () => {
  const [balances, setBalances] = useState<NetworkTokenBalances>({})
  const { address } = getAccount()
  const availableChains = Object.keys(BRIDGABLE_TOKENS)

  useEffect(() => {
    const fetchBalancesAcrossNetworks = async () => {
      const balanceRecord = {}
      availableChains.forEach(async (chainId) => {
        const currentChainId = Number(chainId)
        const currentChainTokens = BRIDGABLE_TOKENS[chainId]
        const tokenBalances: TokenAndBalance[] = await getTokensByChainId(
          address,
          currentChainTokens,
          currentChainId
        )
        balanceRecord[currentChainId] = tokenBalances

        const tokenAllowances = await getTokensAllowance(
          address,
          ROUTER_ADDRESS,
          currentChainTokens,
          currentChainId
        )

        const mergedBalancesAndAllowances = mergeBalancesAndAllowances(
          tokenBalances,
          tokenAllowances
        )
        console.log(
          'mergedBalancesAndAllowances chainId:',
          chainId,
          mergedBalancesAndAllowances
        )
      })
      setBalances(balanceRecord)
    }
    fetchBalancesAcrossNetworks()
  }, [])

  return balances
}

const getTokensAllowance = async (
  owner: string,
  spender: string,
  tokens: Token[],
  chainId: number
): Promise<any> => {
  const inputs = tokens.map((token: Token) => {
    const tokenAddress = token.addresses[
      chainId as keyof Token['addresses']
    ] as `0x${string}`
    return {
      address: tokenAddress,
      abi: erc20ABI,
      functionName: 'allowance',
      chainId,
      args: [owner, spender],
    }
  })
  const allowances: unknown[] = await multicall({
    contracts: inputs,
    chainId,
  })

  return tokens.map((token: Token, index: number) => {
    return {
      token,
      allowance: allowances[index],
    }
  })
}
