// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @title Cofre — Smart Contract de exemplo para o curso Web3
/// @notice Um cofre simples onde usuários depositam e sacam ETH
/// @dev Usado como exemplo didático — NÃO use em produção sem auditoria!

contract Cofre {
    // ═══════════════════════════════════════════
    // VARIÁVEIS DE ESTADO (armazenadas na blockchain)
    // ═══════════════════════════════════════════
    
    address public dono;                          // quem fez deploy
    mapping(address => uint256) public saldos;    // saldo de cada endereço
    uint256 public totalDepositado;               // total geral
    
    // ═══════════════════════════════════════════
    // EVENTOS (logs que Go pode escutar via WebSocket)
    // ═══════════════════════════════════════════
    
    event Deposito(address indexed de, uint256 valor, uint256 timestamp);
    event Saque(address indexed para, uint256 valor, uint256 timestamp);
    event TransferenciaDono(address indexed antigosDono, address indexed novoDono);
    
    // ═══════════════════════════════════════════
    // MODIFIERS (validações reutilizáveis)
    // ═══════════════════════════════════════════
    
    modifier apenasDono() {
        require(msg.sender == dono, "Somente o dono pode fazer isso");
        _;
    }
    
    modifier saldoSuficiente(uint256 valor) {
        require(saldos[msg.sender] >= valor, "Saldo insuficiente");
        _;
    }
    
    // ═══════════════════════════════════════════
    // CONSTRUCTOR (roda UMA VEZ, no deploy)
    // ═══════════════════════════════════════════
    
    constructor() {
        dono = msg.sender;
    }
    
    // ═══════════════════════════════════════════
    // FUNÇÕES
    // ═══════════════════════════════════════════
    
    /// @notice Deposita ETH no cofre
    function depositar() public payable {
        require(msg.value > 0, "Valor deve ser maior que zero");
        
        saldos[msg.sender] += msg.value;
        totalDepositado += msg.value;
        
        emit Deposito(msg.sender, msg.value, block.timestamp);
    }
    
    /// @notice Saca ETH do cofre
    /// @param valor Quantidade de Wei para sacar
    function sacar(uint256 valor) public saldoSuficiente(valor) {
        saldos[msg.sender] -= valor;
        
        // Enviar ETH de volta
        (bool sucesso, ) = payable(msg.sender).call{value: valor}("");
        require(sucesso, "Transferencia falhou");
        
        emit Saque(msg.sender, valor, block.timestamp);
    }
    
    /// @notice Consulta saldo (view = só leitura, não gasta gas)
    /// @param conta Endereço para consultar
    /// @return Saldo em Wei
    function consultarSaldo(address conta) public view returns (uint256) {
        return saldos[conta];
    }
    
    /// @notice Saldo total do contrato
    /// @return Total de ETH no contrato em Wei
    function saldoContrato() public view returns (uint256) {
        return address(this).balance;
    }
    
    /// @notice Transfere a propriedade do cofre
    /// @param novoDono Endereço do novo dono
    function transferirDono(address novoDono) public apenasDono {
        require(novoDono != address(0), "Endereco invalido");
        
        address antigoDono = dono;
        dono = novoDono;
        
        emit TransferenciaDono(antigoDono, novoDono);
    }
    
    /// @notice Permite receber ETH diretamente (sem chamar função)
    receive() external payable {
        saldos[msg.sender] += msg.value;
        totalDepositado += msg.value;
        emit Deposito(msg.sender, msg.value, block.timestamp);
    }
}
