import { useDispatch, useSelector } from 'react-redux'
import { createAsyncThunk } from '@reduxjs/toolkit'
import { getAccount, Address } from '@wagmi/core'

import { AppDispatch, RootState } from '@/store/store'
import { useAppSelector } from '@/store/hooks'
import { FetchState } from './reducer'
import {
  fetchPortfolioBalances,
  NetworkTokenBalancesAndAllowances,
} from '@/utils/actions/fetchPortfolioBalances'
import { getTokenAllowance } from './../../utils/actions/getTokenAllowance'

export const usePortfolioState = (): RootState['portfolio'] => {
  return useAppSelector((state) => state.portfolio)
}

export const usePortfolioBalances = (): NetworkTokenBalancesAndAllowances => {
  return useAppSelector((state) => state.portfolio.balancesAndAllowances)
}

export const fetchAndStoreSingleTokenAllowance = createAsyncThunk(
  'portfolio/fetchAndStoreSingleTokenAllowance',
  async ({
    routerAddress,
    tokenAddress,
    address,
    chainId,
  }: {
    routerAddress: Address
    tokenAddress: Address
    address: Address
    chainId: number
  }) => {
    const allowance = await getTokenAllowance(
      routerAddress,
      tokenAddress,
      address,
      chainId
    )
    return { chainId, tokenAddress, allowance }
  }
)

export const fetchAndStorePortfolioBalances = createAsyncThunk(
  'portfolio/fetchAndStorePortfolioBalances',
  async (address: string) => {
    const portfolioData = await fetchPortfolioBalances(address)
    return portfolioData
  }
)

export const fetchAndStoreSingleNetworkPortfolioBalances = createAsyncThunk(
  'portfolio/fetchAndStoreSingleNetworkPortfolioBalances',
  async ({ address, chainId }: { address: string; chainId: number }) => {
    const portfolioData = await fetchPortfolioBalances(address, chainId)
    return portfolioData
  }
)

export const useFetchPortfolioBalances = (): {
  balancesAndAllowances: NetworkTokenBalancesAndAllowances
  fetchPortfolioBalances: () => void
  status: FetchState
  error: string
} => {
  const dispatch: AppDispatch = useDispatch()
  const { address } = getAccount()
  const { balancesAndAllowances, status, error } = useSelector(
    (state: RootState) => state.portfolio
  )

  const fetch = () => {
    if (address) {
      dispatch(fetchAndStorePortfolioBalances(address))
    }
  }

  return { balancesAndAllowances, fetchPortfolioBalances: fetch, status, error }
}
