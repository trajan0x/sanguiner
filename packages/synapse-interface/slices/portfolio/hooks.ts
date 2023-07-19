import { useDispatch, useSelector } from 'react-redux'
import { getAccount } from '@wagmi/core'

import { AppDispatch, RootState } from '@/store/store'
import { fetchAndStorePortfolioBalances } from './reducer'

export const useFetchPortfolioBalances = () => {
  const dispatch: AppDispatch = useDispatch() // use AppDispatch type instead of generic dispatch
  const { address } = getAccount()
  const { balancesAndAllowances, status, error } = useSelector(
    (state: RootState) => state.portfolio
  )

  const fetch = () => {
    if (address) {
      dispatch(fetchAndStorePortfolioBalances())
    }
  }

  return { balancesAndAllowances, fetchPortfolioBalances: fetch, status, error }
}
