# 📖 Glossário Web3

> Todos os termos que você vai encontrar no curso e no mercado, explicados em português.

---

## A

**ABI (Application Binary Interface)** — "Manual de instruções" de um smart contract. Define quais funções existem, seus parâmetros e retornos. É o que o Go usa para saber como chamar o contrato.

**Address** — Endereço Ethereum. 20 bytes (40 caracteres hex + "0x"). Exemplo: `0x742d35Cc6634C0532925a3b844Bc9e7595f2bD18`. Derivado da chave pública via Keccak-256.

**Airdrop** — Distribuição gratuita de tokens para carteiras. Usado como marketing ou recompensa para early adopters.

**AMM (Automated Market Maker)** — Algoritmo que determina preço em uma DEX usando pools de liquidez e a fórmula x*y=k, em vez de livro de ordens. Uniswap é o mais famoso.

**APY (Annual Percentage Yield)** — Rendimento anual com juros compostos. "Esse pool rende 50% APY" = seu dinheiro cresce 50% ao ano (em teoria).

---

## B

**Block** — Unidade de dados na blockchain. Contém transações, hash do bloco anterior, timestamp e nonce. Como uma "página" do livro contábil.

**Block Explorer** — Ferramenta para visualizar dados da blockchain (blocos, transações, endereços). Etherscan é o mais conhecido. ChainPulse (projeto final) é um mini explorer.

**Bridge** — Protocolo que permite transferir ativos entre blockchains diferentes (ex: Ethereum → Polygon). Ponto comum de hacks.

**Bytecode** — Código compilado de um smart contract que roda na EVM. Humanos escrevem Solidity → compilador gera bytecode.

---

## C

**Chain ID** — Identificador numérico da rede. Mainnet = 1, Sepolia = 11155111, Ganache = 1337. Evita replay attacks entre redes.

**Cold Wallet** — Carteira offline (hardware wallet como Ledger/Trezor). Mais segura porque a chave privada nunca toca a internet.

**Consensus** — Mecanismo pelo qual nós da rede concordam sobre o estado da blockchain. PoW (Bitcoin), PoS (Ethereum), PoA (testnets).

**Contract Account** — Conta controlada por código (smart contract), não por chave privada. Não pode iniciar transações sozinha.

---

## D

**DAO (Decentralized Autonomous Organization)** — Organização governada por smart contracts e votação dos detentores de tokens. Sem CEO — código é a lei.

**dApp (Decentralized Application)** — Aplicação que usa smart contracts como backend. Frontend normal (React/Next.js) + backend na blockchain.

**DeFi (Decentralized Finance)** — Finanças sem intermediários. Lending, trading, seguros — tudo via smart contracts. TVL total: ~$50B+.

**DEX (Decentralized Exchange)** — Exchange sem empresa central. Uniswap, SushiSwap, Curve. Usa AMM para determinar preços.

---

## E

**ECDSA (Elliptic Curve Digital Signature Algorithm)** — Algoritmo criptográfico usado para assinar transações. Sua chave privada assina, qualquer um pode verificar com sua chave pública.

**EIP (Ethereum Improvement Proposal)** — Proposta formal de mudança no Ethereum. EIP-20 = ERC-20, EIP-721 = ERC-721, EIP-1559 = novo modelo de gas.

**EOA (Externally Owned Account)** — Conta controlada por chave privada (sua carteira MetaMask). Pode iniciar transações.

**ERC-20** — Padrão para tokens fungíveis (moedas). Define: `transfer`, `approve`, `transferFrom`, `balanceOf`. USDT, LINK, UNI seguem esse padrão.

**ERC-721** — Padrão para NFTs (tokens únicos). Cada token tem um ID único. Bored Apes, CryptoPunks seguem esse padrão.

**ERC-1155** — Padrão multi-token. Combina fungível + não-fungível no mesmo contrato. Usado em jogos (espadas, poções, skins).

**Etherscan** — O block explorer mais popular do Ethereum. Equivalente do Google para a blockchain.

**EVM (Ethereum Virtual Machine)** — O "processador" que executa smart contracts. Toda rede compatível com EVM (Polygon, BSC, Arbitrum) roda o mesmo código.

---

## F

**Faucet** — Site que distribui ETH gratuito em testnets para desenvolvedores testarem.

**Flash Loan** — Empréstimo instantâneo sem colateral, que deve ser devolvido na MESMA transação. Se não devolver, a transação inteira reverte.

**Fork** — Divisão da blockchain em duas. Soft fork (compatível) vs Hard fork (incompatível). Também: copiar/modificar código de outro projeto.

**Fungível** — Intercambiável. 1 BRL = 1 BRL. 1 USDT = 1 USDT. Oposto de NFT.

---

## G

**Ganache** — Blockchain Ethereum local para desenvolvimento. Cria 10 contas com 1000 ETH cada. Você usa nos exemplos do curso.

**Gas** — Unidade de medida de computação na EVM. Transferência simples = 21.000 gas. Deploy de contrato = 500.000+ gas.

**Gas Limit** — Máximo de gas que você aceita gastar. Se a transação precisar de mais → ela falha (mas você paga o gas consumido até ali).

**Gas Price** — Quanto você paga por unidade de gas, em Gwei. Mais alto = transação processa mais rápido.

**Geth (Go-Ethereum)** — Implementação oficial do Ethereum, escrita em Go. É o client que você usa via `ethclient` no curso.

**Gwei** — Unidade de ETH. 1 ETH = 10⁹ Gwei = 1.000.000.000 Gwei. Usado para gas prices.

---

## H

**Hardhat** — Framework de desenvolvimento Ethereum (alternativa ao Ganache). Rede local + testes + deploy.

**Hash** — Função que transforma dados de qualquer tamanho em um valor fixo. SHA-256 gera 32 bytes. Usado para: IDs de blocos, integridade, mining.

**Hot Wallet** — Carteira online (MetaMask, app do celular). Conveniente, mas menos segura que cold wallet.

---

## I

**Immutable** — Não pode ser alterado depois de criado. Smart contracts são imutáveis na blockchain (a não ser que usem proxy pattern).

**Impermanent Loss (IL)** — Perda temporária que provedores de liquidez sofrem quando o preço dos ativos muda. Chamada "impermanent" porque desaparece se o preço voltar ao original.

**Indexer** — Serviço que lê dados da blockchain e armazena em banco de dados para consulta rápida. The Graph e seu projeto ChainPulse são indexers.

**IPFS (InterPlanetary File System)** — Sistema de armazenamento descentralizado. NFT metadata e imagens ficam no IPFS, não na blockchain (que é cara demais para armazenar arquivos grandes).

---

## K

**Keccak-256** — Função hash usada pelo Ethereum (variante do SHA-3). Gera endereços, topic hashes de eventos, etc.

---

## L

**Layer 1 (L1)** — Blockchain principal (Ethereum, Bitcoin, Solana). Segurança máxima, mas mais lenta e cara.

**Layer 2 (L2)** — Rede que roda "em cima" de uma L1 para ser mais rápida e barata. Arbitrum, Optimism, zkSync. Herda a segurança da L1.

**Liquidity Pool** — Par de tokens depositados em um smart contract (AMM). Provedores ganham taxas de trading.

---

## M

**Mainnet** — Rede principal (dinheiro real). Chain ID = 1 no Ethereum.

**Mempool** — "Sala de espera" de transações que foram enviadas mas ainda não mineradas. Bots de MEV monitoram a mempool.

**Merkle Tree** — Árvore de hashes que permite verificar se um dado faz parte de um conjunto sem baixar tudo. Blocos usam Merkle Trees para transações.

**MEV (Maximal Extractable Value)** — Lucro que mineradores/validadores podem extrair reorganizando, incluindo ou excluindo transações. Front-running, sandwich attacks.

**Mint** — Criar um novo token ou NFT. "Mintar um NFT" = criar e registrar na blockchain.

**Multi-sig** — Carteira que requer múltiplas assinaturas para executar transações. Ex: 3 de 5 sócios precisam assinar para mover fundos.

---

## N

**NFT (Non-Fungible Token)** — Token único e não-intercambiável. Cada um tem um ID diferente. Usado para arte digital, ingressos, identidade.

**Node** — Computador que mantém cópia da blockchain e valida transações. Geth é um node Ethereum.

**Nonce** — Contador de transações de um endereço. Primeira TX = nonce 0, segunda = nonce 1, etc. Também: número que mineradores variam no Proof of Work.

---

## O

**Oracle** — Serviço que traz dados do mundo real para smart contracts. Chainlink é o mais usado. Contratos não podem acessar APIs externas sozinhos.

---

## P

**Private Key** — Número secreto de 256 bits que controla uma conta. NUNCA compartilhe. Perde = perde tudo para sempre.

**Proof of Stake (PoS)** — Consenso onde validadores depositam ETH como garantia. Se agirem mal, perdem o depósito (slashing). Ethereum usa isso desde 2022.

**Proof of Work (PoW)** — Consenso onde mineradores resolvem puzzles computacionais. Bitcoin usa isso. Consome muita energia.

**Proxy Pattern** — Técnica para "atualizar" smart contracts imutáveis. O proxy delega chamadas para um contrato de implementação que pode ser trocado.

**Public Key** — Derivada da private key. Usada para gerar o endereço e verificar assinaturas. Seguro compartilhar.

---

## R

**RPC (Remote Procedure Call)** — Protocolo para comunicar com um nó Ethereum. `ethclient.Dial("http://localhost:8545")` usa JSON-RPC.

**Rug Pull** — Golpe onde criadores de um projeto retiram toda a liquidez e fogem. Os tokens dos investidores viram pó.

---

## S

**Seed Phrase (Mnemonic)** — 12 ou 24 palavras que geram todas as suas chaves privadas. "abandon ability able about..." Perde = perde tudo.

**Slippage** — Diferença entre preço esperado e preço real de um swap no AMM. Quanto maior a ordem relativa ao pool, maior o slippage.

**Smart Contract** — Programa que roda na blockchain. Imutável, transparente, automático. Escrito em Solidity, compilado para bytecode EVM.

**Solidity** — Linguagem mais usada para escrever smart contracts no Ethereum. Sintaxe parecida com JavaScript.

**Staking** — Depositar tokens para ajudar a validar a rede e ganhar recompensas. No Ethereum: 32 ETH mínimo para ser validador.

---

## T

**Testnet** — Rede de teste (ETH sem valor real). Sepolia, Goerli. Para desenvolvedores testarem antes de ir para mainnet.

**Token** — Ativo digital criado por um smart contract. Diferente de ETH (que é nativo da blockchain), tokens são criados por contratos.

**TVL (Total Value Locked)** — Valor total depositado em protocolos DeFi. Métrica principal de saúde do ecossistema.

---

## W

**Wallet** — Software ou hardware que armazena chaves privadas e permite interagir com a blockchain. MetaMask, Ledger, Rainbow.

**Wei** — Menor unidade de ETH. 1 ETH = 10¹⁸ Wei. Como centavos são para o Real, wei é para ETH.

**Web3** — Visão de internet descentralizada baseada em blockchain. Web1 = ler. Web2 = ler + escrever. Web3 = ler + escrever + possuir.

**Whale** — Endereço que possui grande quantidade de tokens. Movimentos de whales afetam o mercado significativamente.

**Wrapped Token** — Token que representa outro ativo em outra blockchain. WETH (Wrapped ETH) é um ERC-20 que representa ETH 1:1.

---

## Y

**Yield** — Rendimento. "Yield farming" = estratégia de maximizar rendimento movendo fundos entre protocolos DeFi.

---

## Números e Símbolos

**0x** — Prefixo que indica valor hexadecimal. Endereços, hashes e dados em Ethereum sempre começam com 0x.

**51% Attack** — Se alguém controla >50% do poder de mineração/staking, pode manipular a blockchain. Na prática é quase impossível em redes grandes.

---

> 💡 Não precisa decorar tudo. Use este glossário como referência enquanto estuda os módulos.
