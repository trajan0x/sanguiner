import _ from 'lodash'
import { useEffect, useState, useMemo } from 'react'
import Slider from 'react-input-slider'
import { stringToBigNum } from '@/utils/stringToBigNum'

import { BigNumber } from '@ethersproject/bignumber'
import { formatUnits } from '@ethersproject/units'
import { useSynapseContext } from '@/utils/providers/SynapseProvider'

import { getCoinTextColorCombined } from '@styles/tokens'
import { calculateExchangeRate } from '@utils/calculateExchangeRate'
import { ALL } from '@constants/withdrawTypes'
import Grid from '@tw/Grid'
import TokenInput from '@components/TokenInput'
import RadioButton from '@components/buttons/RadioButton'
import ReceivedTokenSection from '../components/ReceivedTokenSection'
import PriceImpactDisplay from '../components/PriceImpactDisplay'

import { Transition } from '@headlessui/react'
import { TransactionButton } from '@/components/buttons/TransactionButton'
import { Zero } from '@ethersproject/constants'
import { Token } from '@types'
import { approve, withdraw } from '@/utils/actions/approveAndWithdraw'
import { getTokenAllowance } from '@/utils/actions/getTokenAllowance'
import { PoolData, PoolUserData } from '@types'

const DEFAULT_WITHDRAW_QUOTE = {
  priceImpact: Zero,
  outputs: {},
  allowance: undefined,
  routerAddress: '',
}

const Withdraw = ({
  pool,
  chainId,
  address,
  poolData,
  poolUserData,
  refetchCallback,
}: {
  pool: any
  chainId: number
  address: string
  poolData: PoolData
  poolUserData: PoolUserData
  refetchCallback: () => void
}) => {
  const [inputValue, setInputValue] = useState<{
    bn: BigNumber
    str: string
  }>({ bn: Zero, str: '' })

  const [withdrawQuote, setWithdrawQuote] = useState<{
    priceImpact: BigNumber
    outputs: Record<
      string,
      {
        value: BigNumber
        index: number
      }
    >
    allowance: BigNumber
    routerAddress: string
  }>(DEFAULT_WITHDRAW_QUOTE)

  const [withdrawType, setWithdrawType] = useState(ALL)
  const [percentage, setPercentage] = useState(0)
  const [time, setTime] = useState(Date.now())

  const resetInput = () => {
    setInputValue({ bn: Zero, str: '' })
  }
  const { synapseSDK } = useSynapseContext()

  const sumBigNumbers = (pool: Token, bigNumMap: any) => {
    let sum = Zero
    pool?.poolTokens &&
      pool.poolTokens.map((token) => {
        if (bigNumMap[token.addresses[chainId]]) {
          sum = sum.add(
            bigNumMap[token.addresses[chainId]].value.mul(
              BigNumber.from(10).pow(18 - token.decimals[chainId])
            )
          )
        }
      })
    return sum
  }
  const calculateMaxWithdraw = async () => {
    if (poolUserData == null || address == null) {
      return
    }
    try {
      const outputs: Record<
        string,
        {
          value: BigNumber
          index: number
        }
      > = {}
      if (withdrawType == ALL) {
        const { amounts } = await synapseSDK.calculateRemoveLiquidity(
          chainId,
          pool.swapAddresses[chainId],
          inputValue.bn
        )
        for (const tokenAddr in amounts) {
          outputs[tokenAddr] = amounts[tokenAddr]
        }
      } else {
        const { amount } = await synapseSDK.calculateRemoveLiquidityOne(
          chainId,
          pool.swapAddresses[chainId],
          inputValue.bn,
          withdrawType
        )
        outputs[withdrawType] = amount
      }
      const tokenSum = sumBigNumbers(pool, outputs)
      const priceImpact = calculateExchangeRate(
        inputValue.bn,
        18,
        inputValue.bn.sub(tokenSum),
        18
      )
      const allowance = await getTokenAllowance(
        pool.swapAddresses[chainId],
        pool.addresses[chainId],
        address,
        chainId
      )
      setWithdrawQuote({
        priceImpact,
        allowance,
        outputs,
        routerAddress: pool.swapAddresses[chainId],
      })
    } catch (e) {
      console.log(e)
    }
  }

  useEffect(() => {
    if (poolUserData && poolData && address && pool && inputValue.bn.gt(Zero)) {
      calculateMaxWithdraw()
    }
  }, [inputValue, time, withdrawType])

  const onPercentChange = (percent: number) => {
    if (percent > 100) {
      percent = 100
    }
    setPercentage(percent)
    const numericalOut = poolUserData.lpTokenBalance
      ? formatUnits(
          poolUserData.lpTokenBalance.mul(Number(percent)).div(100),
          pool.decimals[chainId]
        )
      : ''
    onChangeInputValue(pool, numericalOut)
  }

  const onChangeInputValue = (token: Token, value: string) => {
    const bigNum = stringToBigNum(value, token.decimals[chainId])
    if (poolUserData.lpTokenBalance.isZero()) {
      setInputValue({ bn: bigNum, str: value })

      setPercentage(0)
      return
    }
    const pn = bigNum
      ? bigNum.mul(100).div(poolUserData.lpTokenBalance).toNumber()
      : 0
    setInputValue({ bn: bigNum, str: value })

    if (pn > 100) {
      setPercentage(100)
    } else {
      setPercentage(pn)
    }
  }

  let isFromBalanceEnough = true
  let isAllowanceEnough = true

  const getButtonProperties = () => {
    let properties = {
      label: 'Withdraw',
      pendingLabel: 'Withdrawing funds...',
      className: '',
      disabled: false,
      buttonAction: () =>
        withdraw(
          pool,
          'ONE_TENTH',
          null,
          inputValue.bn,
          chainId,
          withdrawType,
          withdrawQuote.outputs
        ),
      postButtonAction: () => {
        refetchCallback()
        setPercentage(0)
        setWithdrawQuote(DEFAULT_WITHDRAW_QUOTE)
        resetInput()
      },
    }

    if (inputValue.bn.eq(0)) {
      properties.label = `Enter amount`
      properties.disabled = true
      return properties
    }

    if (!isFromBalanceEnough) {
      properties.label = `Insufficient Balance`
      properties.disabled = true
      return properties
    }

    if (!isAllowanceEnough) {
      properties.label = `Approve Token(s)`
      properties.pendingLabel = `Approving Token(s)`
      properties.className = 'from-[#feba06] to-[#FEC737]'
      properties.disabled = false
      properties.buttonAction = () =>
        approve(pool, withdrawQuote, inputValue.bn, chainId)
      properties.postButtonAction = () => setTime(0)
      return properties
    }

    return properties
  }

  if (
    withdrawQuote.allowance &&
    !inputValue.bn.isZero() &&
    inputValue.bn.gt(withdrawQuote.allowance)
  ) {
    isAllowanceEnough = false
  }

  if (
    !inputValue.bn.isZero() &&
    inputValue.bn.gt(poolUserData.lpTokenBalance)
  ) {
    isFromBalanceEnough = false
  }

  const {
    label: btnLabel,
    pendingLabel,
    className: btnClassName,
    buttonAction,
    postButtonAction,
    disabled,
  } = useMemo(getButtonProperties, [
    isFromBalanceEnough,
    isAllowanceEnough,
    address,
    inputValue,
    withdrawQuote,
  ])

  const actionBtn = useMemo(
    () => (
      <TransactionButton
        className={btnClassName}
        disabled={disabled}
        onClick={() => buttonAction()}
        onSuccess={() => postButtonAction()}
        label={btnLabel}
        pendingLabel={pendingLabel}
      />
    ),
    [
      buttonAction,
      postButtonAction,
      btnLabel,
      pendingLabel,
      btnClassName,
      isFromBalanceEnough,
      isAllowanceEnough,
    ]
  )

  return (
    <div>
      <div className="percentage">
        <span className="mr-2 text-white">Withdraw Percentage %</span>
        <input
          className={`
            px-2 py-1 w-1/5 rounded-md
            focus:ring-indigo-500 focus:outline-none focus:border-purple-700
            border border-transparent
            bg-[#111111]
            text-gray-300
          `}
          placeholder="0"
          onChange={(e) => {
            onPercentChange(Number(e.currentTarget.value))
          }}
          onFocus={(e) => e.target.select()}
          value={percentage ?? ''}
        />
        <div className="my-2">
          <Slider
            axis="x"
            xstep={10}
            xmin={0}
            xmax={100}
            x={percentage ?? 100}
            onChange={(i) => {
              onPercentChange(i.x)
            }}
            styles={{
              track: {
                backgroundColor: '#E0E7FF',
                width: '95%',
              },
              active: {
                backgroundColor: '#B286FF',
              },
              thumb: {
                backgroundColor: '#CE55FE',
              },
            }}
          />
        </div>
        {/* {error && (
          <div className="text-red-400 opacity-80">{error?.message}</div>
        )} */}
      </div>
      <Grid gap={2} cols={{ xs: 1 }} className="mt-2">
        <RadioButton
          checked={withdrawType === ALL}
          onChange={() => {
            setWithdrawType(ALL)
          }}
          label="Combo"
          labelClassName={withdrawType === ALL && 'text-indigo-500'}
        />
        {pool?.poolTokens &&
          pool.poolTokens.map((token) => {
            const checked = withdrawType === token.addresses[chainId]
            return (
              <RadioButton
                radioClassName={getCoinTextColorCombined(token.color)}
                key={token?.symbol}
                checked={checked}
                onChange={() => {
                  setWithdrawType(token.addresses[chainId])
                }}
                labelClassName={
                  checked &&
                  `${getCoinTextColorCombined(token.color)} opacity-90`
                }
                label={token.name}
              />
            )
          })}
      </Grid>
      <TokenInput
        token={pool}
        key={pool?.symbol}
        inputValueStr={inputValue.str}
        balanceStr={poolUserData?.lpTokenBalanceStr ?? '0.0000'}
        onChange={(value) => onChangeInputValue(pool, value)}
        chainId={chainId}
        address={address}
      />
      {actionBtn}

      <Transition
        appear={true}
        unmount={false}
        show={inputValue.bn.gt(0)}
        enter="transition duration-100 ease-out"
        enterFrom="transform-gpu scale-y-0 "
        enterTo="transform-gpu scale-y-100 opacity-100"
        leave="transition duration-75 ease-out "
        leaveFrom="transform-gpu scale-y-100 opacity-100"
        leaveTo="transform-gpu scale-y-0 "
        className="-mx-6 origin-top "
      >
        <div
          className={`py-3.5 pr-6 pl-6 mt-2 rounded-b-2xl bg-bgBase transition-all`}
        >
          <Grid cols={{ xs: 2 }}>
            <div>
              <ReceivedTokenSection
                poolTokens={pool?.poolTokens ?? []}
                withdrawQuote={withdrawQuote}
                chainId={chainId}
              />
            </div>
            <div>
              {withdrawQuote.priceImpact &&
                withdrawQuote.priceImpact?.gt(Zero) && (
                  <PriceImpactDisplay priceImpact={withdrawQuote.priceImpact} />
                )}
            </div>
          </Grid>
        </div>
      </Transition>
    </div>
  )
}

export default Withdraw
