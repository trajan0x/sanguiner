import { Provider } from '@ethersproject/abstract-provider'
import invariant from 'tiny-invariant'
import { BigNumber } from '@ethersproject/bignumber'
import { BytesLike } from '@ethersproject/bytes'
import { PopulatedTransaction } from 'ethers'
import { AddressZero, Zero } from '@ethersproject/constants'
import { Interface } from '@ethersproject/abi'
import { Contract } from '@ethersproject/contracts'

import {
  handleNativeToken,
  ETH_NATIVE_TOKEN_ADDRESS,
} from './utils/handleNativeToken'
import { BigintIsh } from './constants'
import { SynapseRouter } from './synapseRouter'
import bridgeAbi from './abi/SynapseBridge.json'

type SynapseRouters = {
  [key: number]: SynapseRouter
}

type Query = [string, string, BigNumber, BigNumber, string] & {
  swapAdapter: string
  tokenOut: string
  minAmountOut: BigNumber
  deadline: BigNumber
  rawParams: string
}

type FeeConfig = [number, BigNumber, BigNumber] & {
  bridgeFee: number
  minFee: BigNumber
  maxFee: BigNumber
}

type PoolToken = { isWeth: boolean | undefined; token: string }

class SynapseSDK {
  public synapseRouters: SynapseRouters
  public providers: { [x: number]: Provider }
  public bridgeAbi: Interface = new Interface(bridgeAbi)

  constructor(chainIds: number[], providers: Provider[]) {
    invariant(
      chainIds.length === providers.length,
      `Amount of chains and providers does not equal`
    )
    this.synapseRouters = {}
    this.providers = {}
    for (let i = 0; i < chainIds.length; i++) {
      this.synapseRouters[chainIds[i]] = new SynapseRouter(
        chainIds[i],
        providers[i]
      )
      this.providers[chainIds[i]] = providers[i]
    }
  }

  public async bridgeQuote(
    originChainId: number,
    destChainId: number,
    tokenIn: string,
    tokenOut: string,
    amountIn: BigintIsh,
    deadline?: BigNumber
  ): Promise<{
    feeAmount?: BigNumber | undefined
    feeConfig?: FeeConfig | undefined
    routerAddress?: string | undefined
    maxAmountOut?: BigNumber | undefined
    originQuery?: Query | undefined
    destQuery?: Query | undefined
  }> {
    tokenOut = handleNativeToken(tokenOut)
    tokenIn = handleNativeToken(tokenIn)
    let originQuery
    let destQuery
    const originRouter: SynapseRouter = this.synapseRouters[originChainId]
    const destRouter: SynapseRouter = this.synapseRouters[destChainId]

    // Set deadline
    if (!deadline) {
      const defaultDeadline = Math.floor(Date.now() / 1000) + 10 * 60
      deadline = BigNumber.from(defaultDeadline)
    }

    // Step 0: find connected bridge tokens on destination
    const bridgeTokens =
      await destRouter.routerContract.getConnectedBridgeTokens(tokenOut)

    if (bridgeTokens.length === 0) {
      throw Error('No bridge tokens found for this route')
    }

    const filteredTokens = bridgeTokens.filter(
      (bridgeToken) =>
        bridgeToken.symbol.length !== 0 && bridgeToken.token !== AddressZero
    )

    // Step 1: perform a call to origin SynapseRouter
    const originQueries = await originRouter.routerContract.getOriginAmountOut(
      tokenIn,
      filteredTokens.map((bridgeToken) => bridgeToken.symbol),
      amountIn
    )

    // Step 2: form a list of Destination Requests
    // In practice, there is no need to pass the requests with amountIn = 0, but we will do it for code simplicity
    const requests: { symbol: string; amountIn: BigintIsh }[] = []

    for (let i = 0; i < filteredTokens.length; i++) {
      requests.push({
        symbol: filteredTokens[i].symbol,
        amountIn: originQueries[i].minAmountOut,
      })
    }

    // Step 3: perform a call to destination SynapseRouter
    const destQueries = await destRouter.routerContract.getDestinationAmountOut(
      requests,
      tokenOut
    )
    // Step 4: find the best query (in practice, we could return them all)
    let destInToken
    let maxAmountOut: BigNumber = BigNumber.from(0)
    for (let i = 0; i < destQueries.length; i++) {
      if (destQueries[i].minAmountOut.gt(maxAmountOut)) {
        maxAmountOut = destQueries[i].minAmountOut
        originQuery = originQueries[i]
        destQuery = destQueries[i]
        destInToken = filteredTokens[i].token
      }
    }

    // Get fee data
    let feeAmount
    let feeConfig

    if (originQuery && destInToken) {
      feeAmount = destRouter.routerContract.calculateBridgeFee(
        destInToken,
        originQuery.minAmountOut
      )
      feeConfig = destRouter.routerContract.fee(destInToken)
    }

    if (originQuery && destQuery) {
      originQuery = [...originQuery] as Query
      originQuery[3] = deadline
      originQuery.deadline = deadline
      destQuery = [...destQuery] as Query
      destQuery[3] = deadline
      destQuery.deadline = deadline
    }

    // Router address so allowance handling be set by client
    const routerAddress = originRouter.routerContract.address

    return {
      feeAmount: await feeAmount,
      feeConfig: await feeConfig,
      routerAddress,
      maxAmountOut,
      originQuery,
      destQuery,
    }
  }

  public async bridge(
    to: string,
    originChainId: number,
    destChainId: number,
    token: string,
    amount: BigintIsh,
    originQuery: {
      swapAdapter: string
      tokenOut: string
      minAmountOut: BigintIsh
      deadline: BigintIsh
      rawParams: BytesLike
    },
    destQuery: {
      swapAdapter: string
      tokenOut: string
      minAmountOut: BigintIsh
      deadline: BigintIsh
      rawParams: BytesLike
    }
  ): Promise<PopulatedTransaction> {
    token = handleNativeToken(token)
    const originRouter: SynapseRouter = this.synapseRouters[originChainId]
    return originRouter.routerContract.populateTransaction.bridge(
      to,
      destChainId,
      token,
      amount,
      originQuery,
      destQuery
    )
  }

  // TODO: add gas from bridge
  public async swapQuote(
    chainId: number,
    tokenIn: string,
    tokenOut: string,
    amountIn: BigintIsh,
    deadline?: BigNumber
  ): Promise<{
    routerAddress?: string | undefined
    maxAmountOut?: BigNumber | undefined
    query?: Query | undefined
  }> {
    tokenOut = handleNativeToken(tokenOut)
    tokenIn = handleNativeToken(tokenIn)
    // Set deadline
    if (!deadline) {
      const defaultDeadline = Math.floor(Date.now() / 1000) + 10 * 60
      deadline = BigNumber.from(defaultDeadline)
    }
    const router: SynapseRouter = this.synapseRouters[chainId]

    // Step 0: get the swap quote
    let query = await router.routerContract.getAmountOut(
      tokenIn,
      tokenOut,
      amountIn
    )

    // Router address so allowance handling be set by client
    const routerAddress = router.routerContract.address
    const maxAmountOut = query.minAmountOut

    if (query) {
      query = [...query] as Query
      query[3] = deadline
      query.deadline = deadline
    }
    return {
      routerAddress,
      maxAmountOut,
      query,
    }
  }

  public async swap(
    chainId: number,
    to: string,
    token: string,
    amount: BigintIsh,
    query: {
      swapAdapter: string
      tokenOut: string
      minAmountOut: BigintIsh
      deadline: BigintIsh
      rawParams: BytesLike
    }
  ): Promise<PopulatedTransaction> {
    token = handleNativeToken(token)
    const originRouter: SynapseRouter = this.synapseRouters[chainId]
    return originRouter.routerContract.populateTransaction.swap(
      to,
      token,
      amount,
      query
    )
  }
  public async getBridgeGas(chainId: number): Promise<BigintIsh> {
    const router: SynapseRouter = this.synapseRouters[chainId]
    const bridgeAddress = await router.routerContract.synapseBridge()
    const bridgeContract = new Contract(
      bridgeAddress,
      this.bridgeAbi,
      this.providers[chainId]
    )
    return bridgeContract.chainGasAmount()
  }

  public async getPoolTokens(
    chainId: number,
    poolAddress: string
  ): Promise<PoolToken[]> {
    const router: SynapseRouter = this.synapseRouters[chainId]
    const poolTokens = await router.routerContract.poolTokens(poolAddress)
    return poolTokens.map((token) => {
      return { token: token.token, isWeth: token?.isWeth }
    })
  }

  public async getPoolInfo(
    chainId: number,
    poolAddress: string
  ): Promise<{ tokens: BigNumber | undefined; lpToken: string | undefined }> {
    const router: SynapseRouter = this.synapseRouters[chainId]
    const poolInfo = await router.routerContract.poolInfo(poolAddress)
    return { tokens: poolInfo?.[0], lpToken: poolInfo?.[1] }
  }

  public async getAllPools(chainId: number): Promise<
    {
      poolAddress: string | undefined
      tokens: PoolToken[] | undefined
      lpToken: string | undefined
    }[]
  > {
    const router: SynapseRouter = this.synapseRouters[chainId]
    const pools = await router.routerContract.allPools()
    const res = pools.map((pool) => {
      return {
        poolAddress: pool?.pool,
        tokens: pool?.tokens.map((token) => {
          return { token: token.token, isWeth: token?.isWeth }
        }),
        lpToken: pool?.lpToken,
      }
    })
    return res
  }

  public async calculateAddLiquidity(
    chainId: number,
    poolAddress: string,
    amounts: Record<string, BigNumber>
  ): Promise<{ amount: BigNumber; routerAddress: string }> {
    const router: SynapseRouter = this.synapseRouters[chainId]
    const poolTokens = await router.routerContract.poolTokens(poolAddress)
    const amountArr: BigNumber[] = []
    poolTokens.map((token) => {
      amountArr.push(amounts[token.token] ?? Zero)
    })
    if (amountArr.filter((amount) => !amount.isZero()).length === 0) {
      return { amount: Zero, routerAddress: router.routerContract.address }
    }
    return {
      amount: await router.routerContract.calculateAddLiquidity(
        poolAddress,
        amountArr
      ),
      routerAddress: router.routerContract.address,
    }
  }

  public async calculateRemoveLiquidity(
    chainId: number,
    poolAddress: string,
    amount: BigNumber
  ): Promise<{
    amounts: Record<string, { value: BigNumber; index: number }>
    routerAddress: string
  }> {
    const router: SynapseRouter = this.synapseRouters[chainId]
    const amounts = await router.routerContract.calculateRemoveLiquidity(
      poolAddress,
      amount
    )
    const poolTokens = await router.routerContract.poolTokens(poolAddress)
    const amountsOut: Record<string, { value: BigNumber; index: number }> = {}
    poolTokens.map((token, index) => {
      amountsOut[token.token] = { value: amounts[index], index }
    })
    return {
      amounts: amountsOut,
      routerAddress: router.routerContract.address,
    }
  }

  public async calculateRemoveLiquidityOne(
    chainId: number,
    poolAddress: string,
    amount: BigNumber,
    token: string
  ): Promise<{
    amount: { value: BigNumber; index: number }
    routerAddress: string
  }> {
    const router: SynapseRouter = this.synapseRouters[chainId]

    let poolIndex = 0
    const poolTokens = await router.routerContract.poolTokens(poolAddress)
    poolTokens.map((poolToken, index) => {
      if (poolToken.token === token) {
        poolIndex = index
      }
    })

    const outAmount = await router.routerContract.calculateWithdrawOneToken(
      poolAddress,
      amount,
      poolIndex
    )

    return {
      amount: { value: outAmount, index: poolIndex },
      routerAddress: router.routerContract.address,
    }
  }
}

export { SynapseSDK, ETH_NATIVE_TOKEN_ADDRESS }