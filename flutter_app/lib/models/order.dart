class Order {
  final String orderId;
  final String status;
  final double total;
  final String createdAt;

  Order({
    required this.orderId,
    required this.status,
    required this.total,
    required this.createdAt,
  });

  factory Order.fromJson(Map<String, dynamic> json) {
    return Order(
      orderId: json['order_id'] ?? '',
      status: json['status'] ?? '',
      total: (json['total'] ?? 0).toDouble(),
      createdAt: json['created_at'] ?? '',
    );
  }
}
