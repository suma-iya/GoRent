class Floor {
  final int id;
  final String name;
  final int rent;
  final String createdAt;
  final int? tenant;
  final String? tenantName;
  final String? status;
  final int? notificationId;
  final bool hasPendingAdvancePayment;

  Floor({
    required this.id,
    required this.name,
    required this.rent,
    required this.createdAt,
    this.tenant,
    this.tenantName,
    this.status,
    this.notificationId,
    this.hasPendingAdvancePayment = false,
  });

  factory Floor.fromJson(Map<String, dynamic> json) {
    return Floor(
      id: json['id'],
      name: json['name'],
      rent: json['rent'],
      createdAt: json['created_at'],
      tenant: json['tenant'],
      tenantName: json['tenant_name'],
      status: json['status'],
      notificationId: json['notification_id'],
      hasPendingAdvancePayment: json['has_pending_advance_payment'] ?? false,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'rent': rent,
      'created_at': createdAt,
      'tenant': tenant,
      'tenant_name': tenantName,
      'status': status,
      'notification_id': notificationId,
      'has_pending_advance_payment': hasPendingAdvancePayment,
    };
  }
} 