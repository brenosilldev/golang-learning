// Módulo 20 — Exercícios: Message Queues & Event-Driven
//
// 🧠 Antes de começar:
//   1. Qual a diferença entre at-least-once e exactly-once delivery?
//   2. O que acontece com uma mensagem num consumer group se o consumer travar?
//   3. O que é uma Dead Letter Queue (DLQ) e quando usamos?
//
// NOTA: Este exercício usa NATS simulado em memória (sem dependências externas).
//       Para usar Kafka/NATS real, veja os exemplos no README.

package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// INFRAESTRUTURA: Message Broker em Memória (não modifique)
// =============================================================================

type Message struct {
	ID        string
	Topic     string
	Key       string
	Payload   []byte
	Timestamp time.Time
	Attempts  int
}

type Broker struct {
	mu        sync.RWMutex
	topics    map[string][]Message
	consumers map[string][]chan Message
}

func NewBroker() *Broker {
	return &Broker{
		topics:    make(map[string][]Message),
		consumers: make(map[string][]chan Message),
	}
}

func (b *Broker) Publish(topic string, msg Message) {
	b.mu.Lock()
	defer b.mu.Unlock()
	msg.Topic = topic
	msg.Timestamp = time.Now()
	b.topics[topic] = append(b.topics[topic], msg)
	for _, ch := range b.consumers[topic] {
		select {
		case ch <- msg:
		default:
		}
	}
}

func (b *Broker) Subscribe(topic string) <-chan Message {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan Message, 100)
	b.consumers[topic] = append(b.consumers[topic], ch)
	return ch
}

// =============================================================================
// EXERCÍCIO 1 — Dead Letter Queue (DLQ) Consumer (🟡 Médio)
// =============================================================================
// Implemente um consumer que:
//   - Processa mensagens de um tópico
//   - Se o processamento falhar, recoloca na fila (até maxRetries vezes)
//   - Após maxRetries falhas, move para o tópico DLQ (topic + ".dlq")
//   - Registra métricas: processadas, falhas, movidas para DLQ

type DLQConsumer struct {
	broker     *Broker
	topic      string
	maxRetries int
	// TODO: adicione campos para métricas (processadas, falhas, dlq)
}

type DLQMetrics struct {
	Processed int
	Failed    int
	MovedToDLQ int
}

func NewDLQConsumer(broker *Broker, topic string, maxRetries int) *DLQConsumer {
	// TODO: inicialize
	panic("implemente NewDLQConsumer")
}

func (c *DLQConsumer) Run(ctx context.Context, handler func(Message) error) {
	// TODO:
	// 1. Subscribe no tópico
	// 2. Para cada mensagem, chame handler
	// 3. Se falhar e attempts < maxRetries: republique com attempts+1
	// 4. Se attempts >= maxRetries: publique em topic+".dlq"
	// 5. Pare quando ctx for cancelado
	panic("implemente DLQConsumer.Run")
}

func (c *DLQConsumer) Metrics() DLQMetrics {
	// TODO: retorne as métricas acumuladas
	panic("implemente DLQConsumer.Metrics")
}

// =============================================================================
// EXERCÍCIO 2 — Outbox Pattern (🟡 Médio)
// =============================================================================
// Simule o Outbox Pattern para garantia de entrega.
// Um "banco de dados" local guarda eventos pendentes.
// Um worker os publica no broker e os marca como publicados.
//
// Garanta que mesmo se o processo crashar no meio:
// - Os eventos nunca são perdidos
// - Os eventos podem ser publicados mais de uma vez (at-least-once)

type OutboxEvent struct {
	ID        string
	Topic     string
	Payload   []byte
	Published bool
	CreatedAt time.Time
}

type Outbox struct {
	mu     sync.Mutex
	events []OutboxEvent
}

func NewOutbox() *Outbox {
	return &Outbox{}
}

// SaveEvent salva um evento no outbox (atomicamente com a "transação de negócio")
func (o *Outbox) SaveEvent(id, topic string, payload []byte) {
	// TODO: adicione o evento com Published=false
	panic("implemente SaveEvent")
}

// Pending retorna todos os eventos ainda não publicados
func (o *Outbox) Pending() []OutboxEvent {
	// TODO: filtre e retorne apenas eventos com Published=false
	panic("implemente Pending")
}

// MarkPublished marca um evento como publicado pelo ID
func (o *Outbox) MarkPublished(id string) {
	// TODO: encontre o evento pelo ID e defina Published=true
	panic("implemente MarkPublished")
}

// StartWorker inicia o worker que publica eventos pendentes periodicamente
func (o *Outbox) StartWorker(ctx context.Context, broker *Broker, interval time.Duration) {
	// TODO:
	// A cada `interval`, leia os eventos pendentes e publique no broker
	// Após publicar com sucesso, chame MarkPublished
	// Pare quando ctx for cancelado
	panic("implemente StartWorker")
}

// =============================================================================
// EXERCÍCIO 3 — Saga Pattern: Checkout Distribuído (🔴 Difícil)
// =============================================================================
// Implemente o padrão Saga choreography para um fluxo de checkout:
//
// Fluxo feliz:
//   pagamento.iniciado → [PagamentoService] → estoque.reservado
//   estoque.reservado  → [EstoqueService]   → pagamento.confirmado
//   pagamento.confirmado → [PagamentoService] → pedido.criado
//
// Compensação (rollback):
//   Se qualquer etapa falhar, publique evento de compensação:
//   estoque.falhou → [EstoqueService.compensar] → estoque.liberado
//   pagamento.falhou → [PagamentoService.compensar] → pagamento.estornado

type SagaEvent struct {
	OrderID string
	Amount  float64
	Items   []string
}

// PagamentoService processa pagamentos (e compensações)
type PagamentoService struct {
	broker  *Broker
	failAt  float64 // se amount > failAt, simula falha
}

func NewPagamentoService(broker *Broker, failAt float64) *PagamentoService {
	return &PagamentoService{broker: broker, failAt: failAt}
}

func (s *PagamentoService) Run(ctx context.Context) {
	// TODO: subscribe em "pagamento.iniciado" e "pagamento.compensar"
	// Para "pagamento.iniciado":
	//   - Se amount > failAt: publique "pagamento.falhou"
	//   - Senão: publique "estoque.reservar"
	// Para "pagamento.compensar":
	//   - Publique "pagamento.estornado" (log do estorno)
	panic("implemente PagamentoService.Run")
}

// EstoqueService reserva e libera estoque
type EstoqueService struct {
	broker    *Broker
	reservado map[string]bool
	mu        sync.Mutex
	failOrders map[string]bool // orders que devem falhar
}

func NewEstoqueService(broker *Broker, failOrders ...string) *EstoqueService {
	fail := make(map[string]bool)
	for _, o := range failOrders {
		fail[o] = true
	}
	return &EstoqueService{
		broker:     broker,
		reservado:  make(map[string]bool),
		failOrders: fail,
	}
}

func (s *EstoqueService) Run(ctx context.Context) {
	// TODO: subscribe em "estoque.reservar" e "estoque.liberar"
	// Para "estoque.reservar":
	//   - Se orderID está em failOrders: publique "estoque.falhou" + "pagamento.compensar"
	//   - Senão: reserve o estoque, publique "pagamento.confirmar"
	// Para "estoque.liberar":
	//   - Libere o estoque (delete de reservado), publique "estoque.liberado"
	panic("implemente EstoqueService.Run")
}

// =============================================================================
// MAIN
// =============================================================================

func testDLQ() {
	fmt.Println("\n=== DEAD LETTER QUEUE ===")
	broker := NewBroker()
	consumer := NewDLQConsumer(broker, "pedidos", 3)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	failCount := 0
	go consumer.Run(ctx, func(msg Message) error {
		failCount++
		if failCount <= 6 {
			return errors.New("processamento falhou")
		}
		fmt.Printf("Processou: %s\n", msg.Payload)
		return nil
	})

	// Publica mensagens
	for i := 0; i < 3; i++ {
		broker.Publish("pedidos", Message{
			ID:      fmt.Sprintf("msg-%d", i),
			Payload: []byte(fmt.Sprintf(`{"pedido": %d}`, i)),
		})
	}

	<-ctx.Done()
	m := consumer.Metrics()
	fmt.Printf("Processadas: %d | Falhas: %d | DLQ: %d\n", m.Processed, m.Failed, m.MovedToDLQ)
}

func testOutbox() {
	fmt.Println("\n=== OUTBOX PATTERN ===")
	broker := NewBroker()
	outbox := NewOutbox()

	// Subscriber que observa eventos publicados
	eventCh := broker.Subscribe("pedidos.criados")
	received := 0
	go func() {
		for range eventCh {
			received++
		}
	}()

	// "Transação de negócio": salva pedido E evento no outbox atomicamente
	outbox.SaveEvent("evt-1", "pedidos.criados", []byte(`{"id":"123"}`))
	outbox.SaveEvent("evt-2", "pedidos.criados", []byte(`{"id":"456"}`))

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	outbox.StartWorker(ctx, broker, 50*time.Millisecond)
	<-ctx.Done()

	time.Sleep(50 * time.Millisecond) // aguarda subscriber
	fmt.Printf("Eventos publicados: %d (esperado: 2)\n", received)
	fmt.Printf("Pendentes após publicação: %d (esperado: 0)\n", len(outbox.Pending()))
}

func testSaga() {
	fmt.Println("\n=== SAGA PATTERN ===")
	broker := NewBroker()
	pagamento := NewPagamentoService(broker, 1000.0) // falha se amount > 1000
	estoque := NewEstoqueService(broker, "order-bad") // falha para order-bad

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go pagamento.Run(ctx)
	go estoque.Run(ctx)

	// Assiste todos os tópicos para ver o fluxo
	topicos := []string{
		"pagamento.iniciado", "estoque.reservar", "pagamento.confirmar",
		"pedido.criado", "estoque.falhou", "pagamento.compensar",
		"pagamento.estornado", "estoque.liberar", "estoque.liberado",
	}
	for _, t := range topicos {
		t := t
		go func() {
			for msg := range broker.Subscribe(t) {
				fmt.Printf("[%s] %s\n", t, msg.Payload)
			}
		}()
	}

	// Pedido normal
	broker.Publish("pagamento.iniciado", Message{
		Payload: []byte(`{"order_id":"order-ok","amount":500}`),
	})

	// Pedido com falha de estoque (saga de compensação)
	time.Sleep(100 * time.Millisecond)
	broker.Publish("pagamento.iniciado", Message{
		Payload: []byte(`{"order_id":"order-bad","amount":100}`),
	})

	<-ctx.Done()
}

func main() {
	testDLQ()
	testOutbox()
	testSaga()
	fmt.Println("\n✅ Exercícios do Módulo 20 concluídos!")
}
