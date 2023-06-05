import _ from 'lodash'

import { WETH } from '@constants/tokens/swapMaster'
import { AVWETH, ETH, WETHE } from '@constants/tokens/master'
import { stringToBigNum } from '@/utils/stringToBigNum'
import { getAddress } from '@ethersproject/address'
import TokenInput from '@components/TokenInput'
import PriceImpactDisplay from '../components/PriceImpactDisplay'
import { useSynapseContext } from '@/utils/providers/SynapseProvider'
import { TransactionButton } from '@/components/buttons/TransactionButton'
import { Zero } from '@ethersproject/constants'
import { Token } from '@types'
import { useState, useEffect, useMemo } from 'react'
import { BigNumber } from '@ethersproject/bignumber'
import { calculateExchangeRate } from '@utils/calculateExchangeRate'
import { getTokenAllowance } from '@/utils/actions/getTokenAllowance'
import { approve, deposit } from '@/utils/actions/approveAndDeposit'
import { QUOTE_POLLING_INTERVAL } from '@/constants/bridge' // TODO CHANGE
import { PoolData, PoolUserData } from '@types'
import LoadingTokenInput from '@components/loading/LoadingTokenInput'

const DEFAULT_DEPOSIT_QUOTE = {
  priceImpact: undefined,
  allowances: {},
  routerAddress: '',
}

const Deposit = ({
  pool,
  chainId,
  address,
  poolData,
  poolUserData,
  refetchCallback,
}: {
  pool: Token
  chainId: number
  address: string
  poolData: PoolData
  poolUserData: PoolUserData
  refetchCallback: () => void
}) => {
  // todo store sum in here?
  const [inputValue, setInputValue] = useState<{
    bn: Record<string, BigNumber>
    str: Record<string, string>
  }>({ bn: {}, str: {} })
  const [depositQuote, setDepositQuote] = useState<{
    priceImpact: BigNumber
    allowances: Record<string, BigNumber>
    routerAddress: string
  }>(DEFAULT_DEPOSIT_QUOTE)
  const [time, setTime] = useState(Date.now())
  const { synapseSDK } = useSynapseContext()

  // TODO move this to utils
  const sumBigNumbersFromState = () => {
    let sum = Zero
    pool?.poolTokens &&
      pool.poolTokens.map((token) => {
        if (!token.addresses[chainId]) return
        const tokenAddress = getAddress(token.addresses[chainId])
        if (inputValue.bn[tokenAddress]) {
          sum = sum.add(
            inputValue.bn[getAddress(token.addresses[chainId])].mul(
              BigNumber.from(10).pow(18 - token.decimals[chainId])
            )
          )
        }
      })
    return sum
  }

  const calculateMaxDeposits = async () => {
    try {
      if (poolUserData == null || address == null) {
        return
      }
      let inputSum = sumBigNumbersFromState()
      if (poolData.totalLocked.gt(0) && inputSum.gt(0)) {
        const { amount } = await synapseSDK.calculateAddLiquidity(
          chainId,
          pool.swapAddresses[chainId],
          inputValue.bn
        )

        let allowances: Record<string, BigNumber> = {}
        for (const [key, value] of Object.entries(inputValue.bn)) {
          allowances[key] = await getTokenAllowance(
            pool.swapAddresses[chainId],
            key,
            address,
            chainId
          )
        }

        const priceImpact = calculateExchangeRate(
          inputSum,
          18,
          inputSum.sub(amount),
          18
        )
        // TODO: DOUBLE CHECK THIS
        setDepositQuote({
          priceImpact,
          allowances,
          routerAddress: pool.swapAddresses[chainId],
        })
      } else {
        setDepositQuote(DEFAULT_DEPOSIT_QUOTE)
      }
    } catch (e) {
      console.log(e)
    }
  }

  useEffect(() => {
    const interval = setInterval(
      () => setTime(Date.now()),
      QUOTE_POLLING_INTERVAL
    )
    return () => {
      clearInterval(interval)
    }
  }, [])

  useEffect(() => {
    calculateMaxDeposits()
  }, [inputValue, time, pool, chainId, address])

  const onChangeInputValue = (token: Token, value: string) => {
    const bigNum = stringToBigNum(value, token.decimals[chainId]) ?? Zero
    if (chainId && token) {
      setInputValue({
        bn: {
          ...inputValue.bn,
          [getAddress(token.addresses[chainId])]: bigNum,
        },
        str: {
          ...inputValue.str,
          [getAddress(token.addresses[chainId])]: value,
        },
      })
    }
  }

  useEffect(() => {
    if (poolData && poolUserData && pool && chainId && address) {
      resetInputs()
    }
  }, [poolUserData])

  const resetInputs = () => {
    let initInputValue: {
      bn: Record<string, BigNumber>
      str: Record<string, string>
    } = { bn: {}, str: {} }
    poolUserData.tokens.map((tokenObj, i) => {
      initInputValue.bn[tokenObj.token.addresses[chainId]] = Zero
      initInputValue.str[tokenObj.token.addresses[chainId]] = ''
    })
    setInputValue(initInputValue)
    setDepositQuote(DEFAULT_DEPOSIT_QUOTE)
  }

  let isFromBalanceEnough = true
  let isAllowanceEnough = true

  const getButtonProperties = () => {
    let properties = {
      label: 'Deposit',
      pendingLabel: 'Depositing funds...',
      className: '',
      disabled: false,
      buttonAction: () =>
        deposit(pool, 'ONE_TENTH', null, inputValue.bn, chainId),
      postButtonAction: () => {
        console.log('Post Button Action')
        refetchCallback()
        resetInputs()
      },
    }

    if (sumBigNumbersFromState().eq(0)) {
      properties.disabled = true
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
        approve(pool, depositQuote, inputValue.bn, chainId)
      properties.postButtonAction = () => setTime(0)
      return properties
    }

    return properties
  }

  for (const [tokenAddr, amount] of Object.entries(inputValue.bn)) {
    if (
      Object.keys(depositQuote.allowances).length > 0 &&
      !amount.isZero() &&
      amount.gt(depositQuote.allowances[tokenAddr])
    ) {
      isAllowanceEnough = false
    }
    poolUserData.tokens.map((tokenObj, i) => {
      if (
        tokenObj.token.addresses[chainId] === tokenAddr &&
        amount.gt(tokenObj.balance)
      ) {
        isFromBalanceEnough = false
      }
    })
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
    depositQuote,
  ])

  const actionBtn = useMemo(
    () => (
      <TransactionButton
        className={btnClassName}
        disabled={sumBigNumbersFromState().eq(0) || disabled}
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
    <div className="flex-col">
      <div className="px-2 pt-1 pb-4 bg-bgLight rounded-xl">
        {pool && poolUserData && poolData ? (
          poolUserData.tokens.map((tokenObj, i) => {
            const balanceToken = correctToken(tokenObj.token)
            return (
              <TokenInput
                token={balanceToken}
                key={balanceToken.symbol}
                balanceStr={String(tokenObj.balanceStr)}
                inputValueStr={inputValue.str[balanceToken.addresses[chainId]]}
                onChange={(value) => onChangeInputValue(balanceToken, value)}
                chainId={chainId}
                address={address}
              />
            )
          })
        ) : (
          <>
            <LoadingTokenInput />
            <LoadingTokenInput />
          </>
        )}
      </div>
      {actionBtn}
      {depositQuote.priceImpact && depositQuote.priceImpact?.gt(Zero) && (
        <PriceImpactDisplay priceImpact={depositQuote.priceImpact} />
      )}
    </div>
  )
}
const correctToken = (token: Token) => {
  let balanceToken: Token | undefined
  if (token.symbol == WETH.symbol) {
    balanceToken = ETH
  } else if (token.symbol == AVWETH.symbol) {
    // token = WETHE
    balanceToken = WETHE
  } else {
    balanceToken = token
  }
  return balanceToken
}

export default Deposit
