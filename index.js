import { JsonRpcProvider } from '@ethersproject/providers';

export const allChains = [
  1,
  10,
  56,
  250,
  288,
  1284,
  1285,
  137,
  43144,
  53935,
  42161,
  1313161554,
  1666600000,
  25,
  1088,
  8217,
  2000,
  7700
]

export const rpcProviders = [
  new JsonRpcProvider('https://rpc.ankr.com/eth'),
  new JsonRpcProvider('https://rpc.ankr.com/optimism'),
  new JsonRpcProvider('https://bsc-dataseed1.ninicoin.io/'),
  new JsonRpcProvider('https://rpc.ftm.tools'),
  new JsonRpcProvider('https://replica-oolong.boba.network/'),
  new JsonRpcProvider('https://rpc.api.moonbeam.network'),
  new JsonRpcProvider('https://rpc.api.moonriver.moonbeam.network'),
  new JsonRpcProvider('https://rpc-mainnet.matic.quiknode.pro'),
  new JsonRpcProvider('https://api.avax.network/ext/bc/C/rpc'),
  new JsonRpcProvider('https://subnets.avax.network/defi-kingdoms/dfk-chain/rpc'),
  new JsonRpcProvider('https://arb1.arbitrum.io/rpc'),
  new JsonRpcProvider('https://mainnet.aurora.dev'),
  new JsonRpcProvider('https://harmony-mainnet.chainstacklabs.com'),
  new JsonRpcProvider('https://evm-cronos.crypto.org'),
  new JsonRpcProvider('https://andromeda.metis.io/?owner=1088'),
  new JsonRpcProvider('https://klaytn.blockpi.network/v1/rpc/public'),
  new JsonRpcProvider('https://rpc.ankr.com/dogechain'),
  new JsonRpcProvider('https://mainnode.plexnode.org:8545'),
]
