import _ from 'lodash'
import { RootState } from '@/store/store'
import { useSelector } from 'react-redux'
import { useDispatch } from 'react-redux'
import Select from 'react-select'
import { Token } from '@/utils/types'

import { setToToken } from '@/slices/bridge/reducer'
import { coinSelectStyles } from './styles/coinSelectStyles'
import { useId } from 'react'

const ImageAndCoin = ({ option }: { option: Token }) => {
  const { icon, symbol, routeSymbol } = option
  return (
    <div className="flex items-center space-x-2" key={option.symbol}>
      <img src={icon.src} className="w-6 h-6" />
      <div className="text-xl">{routeSymbol}</div>
    </div>
  )
}

const ToTokenSelect = () => {
  const { toChainId, toToken, toTokens } = useSelector(
    (state: RootState) => state.bridge
  )

  const dispatch = useDispatch()

  const toTokenOptions = toTokens.map((option) => ({
    label: <ImageAndCoin option={option} />,
    value: option,
  }))

  const handleToTokenChange = (selectedOption) => {
    if (selectedOption) {
      dispatch(setToToken(selectedOption.value))
    } else {
      dispatch(setToToken(null))
    }
  }

  const customFilter = (option, searchInput) => {
    if (searchInput) {
      const searchTerm = searchInput.toLowerCase()
      return (
        option.value.symbol.toLowerCase().includes(searchTerm) ||
        option.value.name.toLowerCase().includes(searchTerm) ||
        option.value.addresses[toChainId].toLowerCase().includes(searchTerm)
      )
    }
    return true
  }

  return (
    <Select
      instanceId={useId()}
      styles={coinSelectStyles}
      key={toToken?.symbol}
      options={toTokenOptions}
      filterOption={customFilter}
      onChange={handleToTokenChange}
      isSearchable={true}
      placeholder={<span className="text-xl text-white">Out</span>}
      value={toTokenOptions.find((option) => option.value === toToken)}
    />
  )
}

export default ToTokenSelect
