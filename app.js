import { JsonRpcProvider } from '@ethersproject/providers';
import { Provider } from '@ethersproject/abstract-provider'
import { SynapseSDK } from '@synapsecns/sdk-router';
import { BigNumber } from '@ethersproject/bignumber';
import express from 'express';
import { allChains, rpcProviders } from '../index.js';


//Setting up RPC providers:
const arbitrumProvider = new JsonRpcProvider('https://arb1.arbitrum.io/rpc');
const avalancheProvider = new JsonRpcProvider('https://api.avax.network/ext/bc/C/rpc');

const app = express();
const port = process.env.PORT || 3000;

//Basic hello world
app.get('/', (req, res) => {
  res.send(allChains,rpcProviders)
});

//Setting up arguments
const chainIds = [42161,43114];
const providers = [ arbitrumProvider, avalancheProvider];

//Set up a SynapseSDK Instance
const Synapse = new SynapseSDK(chainIds, providers);


app.get('/swap/:chain/:fromToken/:toToken/:amount', async(req,res) => {
  const chain = req.params.chain;
  //Need logic here that takes in the chain and the token symbol and returns the token address for that chain (for both the to and From tokens) @simon
  const fromToken = req.params.fromToken;
  const toToken = req.params.toToken;
  //Need logic here that takes in the amount and multiplies it by the decimals for that token on its respective chain @simon
  const amount = req.params.amount;

  const resp = await Synapse.swapQuote(chain, fromToken, toToken, BigNumber.from(amount));

  //Hardcoded implementation for testing purposes only
  // const resp = await Synapse.swapQuote(42161, '0xff970a61a04b1ca14834a43f5de4533ebddb5cc8', '0xfd086bc7cd5c481dcc9c85ebe478a1c0b69fcbb9', BigNumber.from(100000000));


  res.json(resp);
});

app.get('/bridge/:fromChain/:toChain/:fromToken/:toToken/:amount', async(req,res) => {
  const fromChain = req.params.fromChain;
  const toChain = req.params.toChain;
  //Need logic here that takes in the chain and the token symbol and returns the token address for that chain (for both the to and From tokens) @simon
  const fromToken = req.params.fromToken;
  const toToken = req.params.toToken;
  //Need logic here that takes in the amount and multiplies it by the decimals for that token on its respective chain @simon
  const amount = req.params.amount;

  const resp = await Synapse.bridgeQuote(fromChain,toChain, fromToken, toToken, BigNumber.from(amount));

  //Hardcoded implementation for testing purposes only
  // const resp = await Synapse.swapQuote(42161, '0xff970a61a04b1ca14834a43f5de4533ebddb5cc8', '0xfd086bc7cd5c481dcc9c85ebe478a1c0b69fcbb9', BigNumber.from(100000000));


  res.json(resp);
});


app.listen(port, () => {
  console.log('Server listening at ${port}')
});


