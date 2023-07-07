import { useMemo, useState } from 'react'
import { useDispatch } from 'react-redux'
import { disconnect, switchNetwork } from '@wagmi/core'
import { setFromChainId } from '@/slices/bridgeSlice'

type PortfolioConnectButton = {
  portfolioChainId: number
  connectedChainId: number
}

export const PortfolioConnectButton = ({
  portfolioChainId,
  connectedChainId,
}: PortfolioConnectButton) => {
  const isCurrentlyConnectedNetwork: boolean = useMemo(() => {
    return portfolioChainId === connectedChainId
  }, [portfolioChainId, connectedChainId])

  return (
    <div data-test-id="portfolio-connect-button">
      {isCurrentlyConnectedNetwork ? (
        <ConnectedButton />
      ) : (
        <ConnectButton chainId={portfolioChainId} />
      )}
    </div>
  )
}

const ConnectedButton = () => {
  const [isDisconnecting, setIsDisconnecting] = useState<boolean>(false)

  const handleDisconnectNetwork = async () => {
    setIsDisconnecting(true)
    try {
      await disconnect()
    } catch (error) {
      error && setIsDisconnecting(false)
    }
  }

  return (
    <button
      data-test-id="connected-button"
      className={`
      flex items-center justify-center
      text-base text-white px-3 py-1 rounded-3xl
      text-center transform-gpu transition-all duration-75
      border border-solid border-transparent
      hover:border-[#3D3D5C]
      `}
      onClick={handleDisconnectNetwork}
    >
      {isDisconnecting ? (
        <div className="flex flex-row text-sm">
          <div
            className={`
            my-auto ml-auto text-transparent w-2 h-2
            border border-red-300 border-solid rounded-full
            `}
          />
          Disconnecting...
        </div>
      ) : (
        <div className="flex flex-row text-sm">
          <div
            className={`
            my-auto ml-auto mr-2 w-2 h-2
            bg-green-500 rounded-full
            `}
          />
          Connected
        </div>
      )}
    </button>
  )
}

const ConnectButton = ({ chainId }: { chainId: number }) => {
  const [isConnecting, setIsConnecting] = useState<boolean>(false)
  const dispatch = useDispatch()

  const handleConnectNetwork = async () => {
    setIsConnecting(true)
    try {
      await switchNetwork({ chainId: chainId }).then((success) => {
        success && dispatch(setFromChainId(chainId))
      })
    } catch (error) {
      error && setIsConnecting(false)
    }
  }

  return (
    <button
      data-test-id="connect-button"
      className={`
      flex items-right justify-center
      text-base text-white px-3 py-1 rounded-3xl
      text-center transform-gpu transition-all duration-75
      border border-solid border-transparent
      hover:border-[#3D3D5C]
      `}
      onClick={handleConnectNetwork}
    >
      {isConnecting ? (
        <div className="flex flex-row text-sm">
          <div
            className={`
            my-auto ml-auto mr-2 text-transparent w-2 h-2
            border border-green-300 border-solid rounded-full
            `}
          />
          Connecting...
        </div>
      ) : (
        <div className="flex flex-row text-sm">
          <div
            className={`
            my-auto ml-auto mr-2 text-transparent w-2 h-2
            border border-indigo-300 border-solid rounded-full
            `}
          />
          Connect
        </div>
      )}
    </button>
  )
}