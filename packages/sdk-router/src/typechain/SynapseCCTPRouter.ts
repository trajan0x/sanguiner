/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */

import {
  ethers,
  EventFilter,
  Signer,
  BigNumber,
  BigNumberish,
  PopulatedTransaction,
  BaseContract,
  ContractTransaction,
  PayableOverrides,
  CallOverrides,
} from "ethers";
import { BytesLike } from "@ethersproject/bytes";
import { Listener, Provider } from "@ethersproject/providers";
import { FunctionFragment, EventFragment, Result } from "@ethersproject/abi";
import type { TypedEventFilter, TypedEvent, TypedListener } from "./common";

interface SynapseCCTPRouterInterface extends ethers.utils.Interface {
  functions: {
    "adapterSwap(address,address,uint256,address,bytes)": FunctionFragment;
    "bridge(address,uint256,address,uint256,(address,address,uint256,uint256,bytes),(address,address,uint256,uint256,bytes))": FunctionFragment;
    "calculateFeeAmount(address,uint256,bool)": FunctionFragment;
    "feeStructures(address)": FunctionFragment;
    "getConnectedBridgeTokens(address)": FunctionFragment;
    "getDestinationAmountOut(tuple[],address)": FunctionFragment;
    "getOriginAmountOut(address,string[],uint256)": FunctionFragment;
    "synapseCCTP()": FunctionFragment;
  };

  encodeFunctionData(
    functionFragment: "adapterSwap",
    values: [string, string, BigNumberish, string, BytesLike]
  ): string;
  encodeFunctionData(
    functionFragment: "bridge",
    values: [
      string,
      BigNumberish,
      string,
      BigNumberish,
      {
        routerAdapter: string;
        tokenOut: string;
        minAmountOut: BigNumberish;
        deadline: BigNumberish;
        rawParams: BytesLike;
      },
      {
        routerAdapter: string;
        tokenOut: string;
        minAmountOut: BigNumberish;
        deadline: BigNumberish;
        rawParams: BytesLike;
      }
    ]
  ): string;
  encodeFunctionData(
    functionFragment: "calculateFeeAmount",
    values: [string, BigNumberish, boolean]
  ): string;
  encodeFunctionData(
    functionFragment: "feeStructures",
    values: [string]
  ): string;
  encodeFunctionData(
    functionFragment: "getConnectedBridgeTokens",
    values: [string]
  ): string;
  encodeFunctionData(
    functionFragment: "getDestinationAmountOut",
    values: [{ symbol: string; amountIn: BigNumberish }[], string]
  ): string;
  encodeFunctionData(
    functionFragment: "getOriginAmountOut",
    values: [string, string[], BigNumberish]
  ): string;
  encodeFunctionData(
    functionFragment: "synapseCCTP",
    values?: undefined
  ): string;

  decodeFunctionResult(
    functionFragment: "adapterSwap",
    data: BytesLike
  ): Result;
  decodeFunctionResult(functionFragment: "bridge", data: BytesLike): Result;
  decodeFunctionResult(
    functionFragment: "calculateFeeAmount",
    data: BytesLike
  ): Result;
  decodeFunctionResult(
    functionFragment: "feeStructures",
    data: BytesLike
  ): Result;
  decodeFunctionResult(
    functionFragment: "getConnectedBridgeTokens",
    data: BytesLike
  ): Result;
  decodeFunctionResult(
    functionFragment: "getDestinationAmountOut",
    data: BytesLike
  ): Result;
  decodeFunctionResult(
    functionFragment: "getOriginAmountOut",
    data: BytesLike
  ): Result;
  decodeFunctionResult(
    functionFragment: "synapseCCTP",
    data: BytesLike
  ): Result;

  events: {};
}

export class SynapseCCTPRouter extends BaseContract {
  connect(signerOrProvider: Signer | Provider | string): this;
  attach(addressOrName: string): this;
  deployed(): Promise<this>;

  listeners<EventArgsArray extends Array<any>, EventArgsObject>(
    eventFilter?: TypedEventFilter<EventArgsArray, EventArgsObject>
  ): Array<TypedListener<EventArgsArray, EventArgsObject>>;
  off<EventArgsArray extends Array<any>, EventArgsObject>(
    eventFilter: TypedEventFilter<EventArgsArray, EventArgsObject>,
    listener: TypedListener<EventArgsArray, EventArgsObject>
  ): this;
  on<EventArgsArray extends Array<any>, EventArgsObject>(
    eventFilter: TypedEventFilter<EventArgsArray, EventArgsObject>,
    listener: TypedListener<EventArgsArray, EventArgsObject>
  ): this;
  once<EventArgsArray extends Array<any>, EventArgsObject>(
    eventFilter: TypedEventFilter<EventArgsArray, EventArgsObject>,
    listener: TypedListener<EventArgsArray, EventArgsObject>
  ): this;
  removeListener<EventArgsArray extends Array<any>, EventArgsObject>(
    eventFilter: TypedEventFilter<EventArgsArray, EventArgsObject>,
    listener: TypedListener<EventArgsArray, EventArgsObject>
  ): this;
  removeAllListeners<EventArgsArray extends Array<any>, EventArgsObject>(
    eventFilter: TypedEventFilter<EventArgsArray, EventArgsObject>
  ): this;

  listeners(eventName?: string): Array<Listener>;
  off(eventName: string, listener: Listener): this;
  on(eventName: string, listener: Listener): this;
  once(eventName: string, listener: Listener): this;
  removeListener(eventName: string, listener: Listener): this;
  removeAllListeners(eventName?: string): this;

  queryFilter<EventArgsArray extends Array<any>, EventArgsObject>(
    event: TypedEventFilter<EventArgsArray, EventArgsObject>,
    fromBlockOrBlockhash?: string | number | undefined,
    toBlock?: string | number | undefined
  ): Promise<Array<TypedEvent<EventArgsArray & EventArgsObject>>>;

  interface: SynapseCCTPRouterInterface;

  functions: {
    adapterSwap(
      recipient: string,
      tokenIn: string,
      amountIn: BigNumberish,
      tokenOut: string,
      rawParams: BytesLike,
      overrides?: PayableOverrides & { from?: string | Promise<string> }
    ): Promise<ContractTransaction>;

    bridge(
      recipient: string,
      chainId: BigNumberish,
      token: string,
      amount: BigNumberish,
      originQuery: {
        routerAdapter: string;
        tokenOut: string;
        minAmountOut: BigNumberish;
        deadline: BigNumberish;
        rawParams: BytesLike;
      },
      destQuery: {
        routerAdapter: string;
        tokenOut: string;
        minAmountOut: BigNumberish;
        deadline: BigNumberish;
        rawParams: BytesLike;
      },
      overrides?: PayableOverrides & { from?: string | Promise<string> }
    ): Promise<ContractTransaction>;

    calculateFeeAmount(
      token: string,
      amount: BigNumberish,
      isSwap: boolean,
      overrides?: CallOverrides
    ): Promise<[BigNumber] & { fee: BigNumber }>;

    feeStructures(
      token: string,
      overrides?: CallOverrides
    ): Promise<
      [number, BigNumber, BigNumber, BigNumber] & {
        relayerFee: number;
        minBaseFee: BigNumber;
        minSwapFee: BigNumber;
        maxFee: BigNumber;
      }
    >;

    getConnectedBridgeTokens(
      tokenOut: string,
      overrides?: CallOverrides
    ): Promise<
      [([string, string] & { symbol: string; token: string })[]] & {
        tokens: ([string, string] & { symbol: string; token: string })[];
      }
    >;

    getDestinationAmountOut(
      requests: { symbol: string; amountIn: BigNumberish }[],
      tokenOut: string,
      overrides?: CallOverrides
    ): Promise<
      [
        ([string, string, BigNumber, BigNumber, string] & {
          routerAdapter: string;
          tokenOut: string;
          minAmountOut: BigNumber;
          deadline: BigNumber;
          rawParams: string;
        })[]
      ] & {
        destQueries: ([string, string, BigNumber, BigNumber, string] & {
          routerAdapter: string;
          tokenOut: string;
          minAmountOut: BigNumber;
          deadline: BigNumber;
          rawParams: string;
        })[];
      }
    >;

    getOriginAmountOut(
      tokenIn: string,
      tokenSymbols: string[],
      amountIn: BigNumberish,
      overrides?: CallOverrides
    ): Promise<
      [
        ([string, string, BigNumber, BigNumber, string] & {
          routerAdapter: string;
          tokenOut: string;
          minAmountOut: BigNumber;
          deadline: BigNumber;
          rawParams: string;
        })[]
      ] & {
        originQueries: ([string, string, BigNumber, BigNumber, string] & {
          routerAdapter: string;
          tokenOut: string;
          minAmountOut: BigNumber;
          deadline: BigNumber;
          rawParams: string;
        })[];
      }
    >;

    synapseCCTP(overrides?: CallOverrides): Promise<[string]>;
  };

  adapterSwap(
    recipient: string,
    tokenIn: string,
    amountIn: BigNumberish,
    tokenOut: string,
    rawParams: BytesLike,
    overrides?: PayableOverrides & { from?: string | Promise<string> }
  ): Promise<ContractTransaction>;

  bridge(
    recipient: string,
    chainId: BigNumberish,
    token: string,
    amount: BigNumberish,
    originQuery: {
      routerAdapter: string;
      tokenOut: string;
      minAmountOut: BigNumberish;
      deadline: BigNumberish;
      rawParams: BytesLike;
    },
    destQuery: {
      routerAdapter: string;
      tokenOut: string;
      minAmountOut: BigNumberish;
      deadline: BigNumberish;
      rawParams: BytesLike;
    },
    overrides?: PayableOverrides & { from?: string | Promise<string> }
  ): Promise<ContractTransaction>;

  calculateFeeAmount(
    token: string,
    amount: BigNumberish,
    isSwap: boolean,
    overrides?: CallOverrides
  ): Promise<BigNumber>;

  feeStructures(
    token: string,
    overrides?: CallOverrides
  ): Promise<
    [number, BigNumber, BigNumber, BigNumber] & {
      relayerFee: number;
      minBaseFee: BigNumber;
      minSwapFee: BigNumber;
      maxFee: BigNumber;
    }
  >;

  getConnectedBridgeTokens(
    tokenOut: string,
    overrides?: CallOverrides
  ): Promise<([string, string] & { symbol: string; token: string })[]>;

  getDestinationAmountOut(
    requests: { symbol: string; amountIn: BigNumberish }[],
    tokenOut: string,
    overrides?: CallOverrides
  ): Promise<
    ([string, string, BigNumber, BigNumber, string] & {
      routerAdapter: string;
      tokenOut: string;
      minAmountOut: BigNumber;
      deadline: BigNumber;
      rawParams: string;
    })[]
  >;

  getOriginAmountOut(
    tokenIn: string,
    tokenSymbols: string[],
    amountIn: BigNumberish,
    overrides?: CallOverrides
  ): Promise<
    ([string, string, BigNumber, BigNumber, string] & {
      routerAdapter: string;
      tokenOut: string;
      minAmountOut: BigNumber;
      deadline: BigNumber;
      rawParams: string;
    })[]
  >;

  synapseCCTP(overrides?: CallOverrides): Promise<string>;

  callStatic: {
    adapterSwap(
      recipient: string,
      tokenIn: string,
      amountIn: BigNumberish,
      tokenOut: string,
      rawParams: BytesLike,
      overrides?: CallOverrides
    ): Promise<BigNumber>;

    bridge(
      recipient: string,
      chainId: BigNumberish,
      token: string,
      amount: BigNumberish,
      originQuery: {
        routerAdapter: string;
        tokenOut: string;
        minAmountOut: BigNumberish;
        deadline: BigNumberish;
        rawParams: BytesLike;
      },
      destQuery: {
        routerAdapter: string;
        tokenOut: string;
        minAmountOut: BigNumberish;
        deadline: BigNumberish;
        rawParams: BytesLike;
      },
      overrides?: CallOverrides
    ): Promise<void>;

    calculateFeeAmount(
      token: string,
      amount: BigNumberish,
      isSwap: boolean,
      overrides?: CallOverrides
    ): Promise<BigNumber>;

    feeStructures(
      token: string,
      overrides?: CallOverrides
    ): Promise<
      [number, BigNumber, BigNumber, BigNumber] & {
        relayerFee: number;
        minBaseFee: BigNumber;
        minSwapFee: BigNumber;
        maxFee: BigNumber;
      }
    >;

    getConnectedBridgeTokens(
      tokenOut: string,
      overrides?: CallOverrides
    ): Promise<([string, string] & { symbol: string; token: string })[]>;

    getDestinationAmountOut(
      requests: { symbol: string; amountIn: BigNumberish }[],
      tokenOut: string,
      overrides?: CallOverrides
    ): Promise<
      ([string, string, BigNumber, BigNumber, string] & {
        routerAdapter: string;
        tokenOut: string;
        minAmountOut: BigNumber;
        deadline: BigNumber;
        rawParams: string;
      })[]
    >;

    getOriginAmountOut(
      tokenIn: string,
      tokenSymbols: string[],
      amountIn: BigNumberish,
      overrides?: CallOverrides
    ): Promise<
      ([string, string, BigNumber, BigNumber, string] & {
        routerAdapter: string;
        tokenOut: string;
        minAmountOut: BigNumber;
        deadline: BigNumber;
        rawParams: string;
      })[]
    >;

    synapseCCTP(overrides?: CallOverrides): Promise<string>;
  };

  filters: {};

  estimateGas: {
    adapterSwap(
      recipient: string,
      tokenIn: string,
      amountIn: BigNumberish,
      tokenOut: string,
      rawParams: BytesLike,
      overrides?: PayableOverrides & { from?: string | Promise<string> }
    ): Promise<BigNumber>;

    bridge(
      recipient: string,
      chainId: BigNumberish,
      token: string,
      amount: BigNumberish,
      originQuery: {
        routerAdapter: string;
        tokenOut: string;
        minAmountOut: BigNumberish;
        deadline: BigNumberish;
        rawParams: BytesLike;
      },
      destQuery: {
        routerAdapter: string;
        tokenOut: string;
        minAmountOut: BigNumberish;
        deadline: BigNumberish;
        rawParams: BytesLike;
      },
      overrides?: PayableOverrides & { from?: string | Promise<string> }
    ): Promise<BigNumber>;

    calculateFeeAmount(
      token: string,
      amount: BigNumberish,
      isSwap: boolean,
      overrides?: CallOverrides
    ): Promise<BigNumber>;

    feeStructures(token: string, overrides?: CallOverrides): Promise<BigNumber>;

    getConnectedBridgeTokens(
      tokenOut: string,
      overrides?: CallOverrides
    ): Promise<BigNumber>;

    getDestinationAmountOut(
      requests: { symbol: string; amountIn: BigNumberish }[],
      tokenOut: string,
      overrides?: CallOverrides
    ): Promise<BigNumber>;

    getOriginAmountOut(
      tokenIn: string,
      tokenSymbols: string[],
      amountIn: BigNumberish,
      overrides?: CallOverrides
    ): Promise<BigNumber>;

    synapseCCTP(overrides?: CallOverrides): Promise<BigNumber>;
  };

  populateTransaction: {
    adapterSwap(
      recipient: string,
      tokenIn: string,
      amountIn: BigNumberish,
      tokenOut: string,
      rawParams: BytesLike,
      overrides?: PayableOverrides & { from?: string | Promise<string> }
    ): Promise<PopulatedTransaction>;

    bridge(
      recipient: string,
      chainId: BigNumberish,
      token: string,
      amount: BigNumberish,
      originQuery: {
        routerAdapter: string;
        tokenOut: string;
        minAmountOut: BigNumberish;
        deadline: BigNumberish;
        rawParams: BytesLike;
      },
      destQuery: {
        routerAdapter: string;
        tokenOut: string;
        minAmountOut: BigNumberish;
        deadline: BigNumberish;
        rawParams: BytesLike;
      },
      overrides?: PayableOverrides & { from?: string | Promise<string> }
    ): Promise<PopulatedTransaction>;

    calculateFeeAmount(
      token: string,
      amount: BigNumberish,
      isSwap: boolean,
      overrides?: CallOverrides
    ): Promise<PopulatedTransaction>;

    feeStructures(
      token: string,
      overrides?: CallOverrides
    ): Promise<PopulatedTransaction>;

    getConnectedBridgeTokens(
      tokenOut: string,
      overrides?: CallOverrides
    ): Promise<PopulatedTransaction>;

    getDestinationAmountOut(
      requests: { symbol: string; amountIn: BigNumberish }[],
      tokenOut: string,
      overrides?: CallOverrides
    ): Promise<PopulatedTransaction>;

    getOriginAmountOut(
      tokenIn: string,
      tokenSymbols: string[],
      amountIn: BigNumberish,
      overrides?: CallOverrides
    ): Promise<PopulatedTransaction>;

    synapseCCTP(overrides?: CallOverrides): Promise<PopulatedTransaction>;
  };
}
