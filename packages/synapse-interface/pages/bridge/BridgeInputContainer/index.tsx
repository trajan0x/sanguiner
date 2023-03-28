import _ from 'lodash'
import { Zero } from '@ethersproject/constants'
import { formatBNToString } from '@bignumber/format'

import SelectTokenDropdown from './SelectTokenDropdown'
import { ChainLabel } from './ChainLabel'
import { Token } from '@utils/classes/Token'
import { useAccount, useBalance, useNetwork } from 'wagmi'

import { ChainId } from '@constants/networks'

import SwitchButton from '@components/buttons/SwitchButton'
import MiniMaxButton from '@components/buttons/MiniMaxButton'
import { BigNumber } from '@ethersproject/bignumber'
// import { useTokenBalance } from '@hooks/tokens/useTokenBalances'
// import { useNetworkController } from '@hooks/wallet/useNetworkController'

import { cleanNumberInput } from '@utils/cleanNumberInput'

export default function BridgeInputContainer({
  selected,
  inputValue,
  isSwapFrom,
  onChangeAmount,
  swapFromToChains,
  setDisplayType,
  chainId,
  onChangeChain,
  possibleChains,
  connectedChainId,
}: {
  isSwapFrom: boolean
  selected: Token
  onChangeAmount: (v: string) => void
  inputValue: string
  swapFromToChains?: () => void
  tokens: any[]
  chainId: number
  setDisplayType: (v: string) => void
  onChangeChain: (v: number) => void
  possibleChains: string[]
  connectedChainId: number
}) {
  const { address } = useAccount()
  const tokenAddr = selected.addresses[chainId as keyof Token['addresses']]

  let tokenBalance: BigNumber
  if (!tokenAddr) {
    const { data: rawTokenBalance } = useBalance({
      chainId: chainId,
      address: address,
    })
    tokenBalance = rawTokenBalance?.value ?? Zero
  } else {
    const { data: rawTokenBalance } = useBalance({
      chainId: chainId,
      address: address,
      token: `0x${tokenAddr.slice(2)}`,
    })
    tokenBalance = rawTokenBalance?.value ?? Zero
  }

  const formattedBalance = formatBNToString(
    tokenBalance,
    selected.decimals[chainId as keyof Token['decimals']],
    4
  )

  let isConnected = address !== null

  function onChange(e: any) {
    let val = e.target.value

    let cleanValue = cleanNumberInput(val)

    onChangeAmount(cleanValue)
  }

  function onClickBalance() {
    onChangeAmount(
      formatBNToString(
        tokenBalance,
        selected.decimals[chainId as keyof Token['decimals']],
        4
      )
    )
  }

  let balanceLabel
  if (isSwapFrom) {
    balanceLabel = (
      <a onClick={onClickBalance} className="hover:underline group">
        <small className="text-xs text-gray-500 cursor-pointer group-hover:underline">
          Max:{' '}
          <span className="font-medium text-gray-400 ">{formattedBalance}</span>{' '}
          {selected.symbol}
        </small>
      </a>
    )
  }
  return (
    <div
      className={`
        text-left px-2 sm:px-4 pt-2 pb-4 rounded-xl
        bg-bgLight
      `}
    >
      <div>
        <div className="pt-1 pb-3">
          {!isSwapFrom && (
            <div className="absolute">
              <div className="-mt-12">
                <SwitchButton onClick={swapFromToChains ?? (() => null)} />
              </div>
            </div>
          )}
          <ChainLabel
            isSwapFrom={isSwapFrom}
            chainId={chainId}
            setDisplayType={setDisplayType}
            onChangeChain={onChangeChain}
            possibleChains={possibleChains}
            connectedChainId={connectedChainId}
          />
        </div>
      </div>
      <div className="flex h-16 mb-4 space-x-2">
        <div
          className={`
            flex flex-grow items-center
            pl-3 sm:pl-4
            w-full h-20
            rounded-xl
            border border-white border-opacity-20
            ${
              isSwapFrom &&
              ' transform-gpu transition-all duration-75 hover:border-opacity-30'
            }
          `}
        >
          <SelectTokenDropdown
            chainId={chainId}
            selected={selected}
            onClick={() => {
              setDisplayType(isSwapFrom ? 'from' : 'to')
            }}
          />

          <input
            pattern="[0-9.]+"
            disabled={!isSwapFrom} // may cause issues idk goal is to prevent to result from being selectable
            className={`
              ml-4
              ${isSwapFrom && isConnected ? '-mt-0 md:-mt-4' : '-mt-0'}
              focus:outline-none
              bg-transparent
              pr-4
              w-2/3
             placeholder:text-[#88818C]
             text-white text-opacity-80 text-lg md:text-2xl lg:text-2xl font-medium
            `}
            placeholder="0.0000"
            onChange={isSwapFrom ? onChange : () => {}}
            value={inputValue}
            name="inputRow"
          />
          {isSwapFrom && isConnected && (
            <label
              htmlFor="inputRow"
              className="absolute hidden pt-1 mt-8 ml-40 text-xs text-white transition-all duration-150 md:block transform-gpu hover:text-opacity-70 hover:cursor-pointer"
              onClick={onClickBalance}
            >
              {formattedBalance}
              <span className="text-opacity-50 text-secondaryTextColor">
                {' '}
                available
              </span>
            </label>
          )}
          {isSwapFrom && isConnected && (
            <div className="hidden mr-2 sm:inline-block">
              <MiniMaxButton onClickBalance={onClickBalance} />
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
