import { useMemo } from 'react'
import { useNetwork } from 'wagmi'
import { switchNetwork } from '@wagmi/core'
import { Zero } from '@ethersproject/constants'
import { TokenWithBalanceAndAllowance } from '@/utils/hooks/usePortfolioBalances'
import { Chain } from '@/utils/types'
import { CHAINS_BY_ID } from '@/constants/chains'
import Image from 'next/image'
import { Token } from '@/utils/types'
import { BigNumber } from 'ethers'
import { formatBNToString } from '@/utils/bignumber/format'

type SingleNetworkPortfolioProps = {
  chainId: number
  tokens: TokenWithBalanceAndAllowance[]
}

function separateTokensByAllowance(
  tokens: TokenWithBalanceAndAllowance[]
): [TokenWithBalanceAndAllowance[], TokenWithBalanceAndAllowance[]] {
  const tokensWithAllowance: TokenWithBalanceAndAllowance[] = []
  const tokensWithoutAllowance: TokenWithBalanceAndAllowance[] = []

  tokens.forEach((token) => {
    // allowance is null for native gas tokens
    if (token.allowance === null) {
      tokensWithAllowance.push(token)
    } else if (token.allowance.gt(Zero)) {
      tokensWithAllowance.push(token)
    } else {
      tokensWithoutAllowance.push(token)
    }
  })

  return [tokensWithAllowance, tokensWithoutAllowance]
}

function sortByBalanceDescending(
  tokens: TokenWithBalanceAndAllowance[]
): TokenWithBalanceAndAllowance[] {
  return tokens.sort(
    (a: TokenWithBalanceAndAllowance, b: TokenWithBalanceAndAllowance) =>
      b.balance.gt(a.balance) ? 1 : -1
  )
}

export const SingleNetworkPortfolio = ({
  chainId,
  tokens,
}: SingleNetworkPortfolioProps) => {
  const currentChain: Chain = CHAINS_BY_ID[chainId]

  const [tokensWithAllowance, tokensWithoutAllowance] =
    separateTokensByAllowance(tokens)

  const sortedTokensWithAllowance: TokenWithBalanceAndAllowance[] =
    sortByBalanceDescending(tokensWithAllowance)
  const sortedTokensWithoutAllowance: TokenWithBalanceAndAllowance[] =
    sortByBalanceDescending(tokensWithoutAllowance)

  console.log('sortedTokensWithAllowance:', sortedTokensWithAllowance)
  console.log('sortedTokensWithoutAllowance:', sortedTokensWithoutAllowance)

  return (
    <div className="flex flex-col">
      <PortfolioNetwork
        displayName={currentChain.name}
        chainIcon={currentChain.chainImg}
        chainId={chainId}
      />
      {sortedTokensWithAllowance.map(
        ({ token, balance, allowance }: TokenWithBalanceAndAllowance) => (
          <PortfolioTokenAsset
            token={token}
            balance={balance}
            chainId={chainId}
          />
        )
      )}
      <div className="bg-gray-500">
        {sortedTokensWithoutAllowance.map(
          ({ token, balance, allowance }: TokenWithBalanceAndAllowance) => (
            <PortfolioTokenAsset
              token={token}
              balance={balance}
              chainId={chainId}
            />
          )
        )}
      </div>
    </div>
  )
}

type PortfolioTokenAssetProps = {
  token: Token
  balance: BigNumber
  chainId: number
}

const PortfolioTokenAsset = ({
  token,
  balance,
  chainId,
}: PortfolioTokenAssetProps) => {
  const { icon, symbol, decimals, addresses } = token
  const parsedBalance = formatBNToString(balance, decimals[chainId], 3)

  return (
    <div className="flex flex-row my-2 text-white">
      <div className="flex flex-row w-1/2 text-left">
        <Image
          alt={`${symbol} img`}
          className="w-6 h-6 mr-2 rounded-md"
          src={icon}
        />
        <div>{symbol}</div>
      </div>
      <div className="flex flex-row w-1/2 text-left">
        <div>{parsedBalance}</div>
      </div>
    </div>
  )
}

type PortfolioNetworkProps = {
  displayName: string
  chainIcon: string
  chainId: number
}

const PortfolioNetwork = ({
  displayName,
  chainIcon,
  chainId,
}: PortfolioNetworkProps) => {
  const { chain } = useNetwork()
  const isCurrentlyConnectedNetwork: boolean = useMemo(() => {
    return chainId === chain.id
  }, [chain.id])

  return (
    <div className="flex flex-row justify-between">
      <div className="flex flex-row">
        <Image
          className="mr-4 rounded-md w-7 h-7"
          alt={`${displayName} img`}
          src={chainIcon}
        />
        <div className="font-medium text-white text-18">{displayName}</div>
      </div>
      {isCurrentlyConnectedNetwork ? (
        <ConnectedButton />
      ) : (
        <ConnectButton chainId={chainId} />
      )}
    </div>
  )
}

const ConnectedButton = () => {
  const buttonClassName = `
    h-8 flex items-center
    text-base text-white px-4 py-2 rounded-3xl
    text-center transform-gpu transition-all duration-75
    border-2 border-[#D747FF] radial-gradient-bg
    hover:cursor-default
  `

  return <button className={buttonClassName}>Connected</button>
}

const ConnectButton = ({ chainId }: { chainId: number }) => {
  const handleConnectNetwork = async () => {
    await switchNetwork({ chainId: chainId })
  }

  const buttonClassName = `
    h-8 flex items-center
    text-base text-white px-4 py-2 rounded-3xl
    text-center transform-gpu transition-all duration-75
    border-2 border-[#101018]
    hover:cursor-pointer
  `

  return (
    <button className={buttonClassName} onClick={handleConnectNetwork}>
      Connect
    </button>
  )
}

export const PortfolioAssetHeader = () => {
  return (
    <div className="flex text-[#CCCAD3BF] my-2">
      <div className="w-1/2 text-left">Token</div>
      <div className="w-1/2 text-left">Amount</div>
    </div>
  )
}
