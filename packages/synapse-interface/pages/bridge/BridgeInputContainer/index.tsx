import _ from 'lodash'
import { Zero } from '@ethersproject/constants'
import { formatBNToString } from '@bignumber/format'
import { Token } from '@/utils/types'
import SwitchButton from '@components/buttons/SwitchButton'
import MiniMaxButton from '@components/buttons/MiniMaxButton'
import { BigNumber } from '@ethersproject/bignumber'
import { cleanNumberInput } from '@utils/cleanNumberInput'
import { ChainLabel } from './ChainLabel'
import SelectTokenDropdown from './SelectTokenDropdown'

export default function BridgeInputContainer({
  address,
  isOrigin,
  chains,
  chainId,
  inputString,
  selectedToken,
  connectedChainId,
  onChangeChain,
  onChangeAmount,
  setDisplayType,
  handleChainFlip,
  fromTokenBalance,
}: {
  address: `0x${string}`
  isOrigin: boolean
  chains: string[]
  chainId: number
  inputString: string
  selectedToken: Token
  connectedChainId: number
  setDisplayType: (v: string) => void
  handleChainFlip?: () => void
  onChangeAmount?: (v: string) => void
  onChangeChain: (v: number) => void
  fromTokenBalance?: BigNumber
}) {
  let formattedBalance = ''
  if (fromTokenBalance) {
    formattedBalance = formatBNToString(
      fromTokenBalance,
      selectedToken.decimals[chainId as keyof Token['decimals']],
      4
    )
  }
  console.log(fromTokenBalance, `fromTokenBalance`)
  const isConnected = address !== null

  function onChange(e: any) {
    const val = e.target.value

    const cleanValue = cleanNumberInput(val)

    onChangeAmount(cleanValue)
  }

  function onClickBalance() {
    onChangeAmount(
      formatBNToString(
        fromTokenBalance,
        selectedToken.decimals[chainId as keyof Token['decimals']],
        4
      )
    )
  }

  let balanceLabel
  if (isOrigin) {
    balanceLabel = (
      <a onClick={onClickBalance} className="hover:underline group">
        <small className="text-xs text-gray-500 cursor-pointer group-hover:underline">
          Max:{' '}
          <span className="font-medium text-gray-400 ">{formattedBalance}</span>{' '}
          {selectedToken.symbol}
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
          {!isOrigin && (
            <div className="absolute">
              <div className="-mt-12">
                <SwitchButton onClick={handleChainFlip ?? (() => null)} />
              </div>
            </div>
          )}
          <ChainLabel
            isOrigin={isOrigin}
            chainId={chainId}
            setDisplayType={setDisplayType}
            onChangeChain={onChangeChain}
            chains={chains}
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
              isOrigin &&
              ' transform-gpu transition-all duration-75 hover:border-opacity-30'
            }
          `}
        >
          <SelectTokenDropdown
            chainId={chainId}
            selectedToken={selectedToken}
            onClick={() => {
              setDisplayType(isOrigin ? 'from' : 'to')
            }}
          />

          <input
            pattern="[0-9.]+"
            disabled={!isOrigin} // may cause issues idk goal is to prevent to result from being selectable
            className={`
              ml-4
              ${isOrigin && isConnected ? '-mt-0 md:-mt-4' : '-mt-0'}
              focus:outline-none
              bg-transparent
              pr-4
              w-2/3
             placeholder:text-[#88818C]
             text-white text-opacity-80 text-lg md:text-2xl lg:text-2xl font-medium
            `}
            placeholder="0.0000"
            onChange={onChange ?? (() => null)}
            value={inputString}
            name="inputRow"
          />
          {isOrigin && isConnected && (
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
          {isOrigin && isConnected && (
            <div className="hidden mr-2 sm:inline-block">
              <MiniMaxButton onClickBalance={onClickBalance} />
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
