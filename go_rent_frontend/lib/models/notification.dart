class AppNotification {
  final int id;
  final String message;
  final String status;
  final String createdAt;
  final NotificationProperty property;
  final NotificationFloor floor;
  final bool showActions;
  final bool isRead;
  final String? comment;
  final int? senderId;
  final int? receiverId;
  final String? senderName;
  final String? receiverName;

  AppNotification({
    required this.id,
    required this.message,
    required this.status,
    required this.createdAt,
    required this.property,
    required this.floor,
    required this.showActions,
    required this.isRead,
    this.comment,
    this.senderId,
    this.receiverId,
    this.senderName,
    this.receiverName,
  });

  factory AppNotification.fromJson(Map<String, dynamic> json) {
    return AppNotification(
      id: json['id'],
      message: json['message'],
      status: json['status'],
      createdAt: json['created_at'],
      property: NotificationProperty.fromJson(json['property']),
      floor: NotificationFloor.fromJson(json['floor']),
      showActions: json['show_actions'] ?? false,
      isRead: json['is_read'] ?? false,
      comment: json['comment'],
      senderId: json['sender_id'],
      receiverId: json['receiver_id'],
      senderName: json['sender_name'],
      receiverName: json['receiver_name'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'message': message,
      'status': status,
      'created_at': createdAt,
      'property': property.toJson(),
      'floor': floor.toJson(),
      'show_actions': showActions,
      'is_read': isRead,
      'comment': comment,
      'sender_id': senderId,
      'receiver_id': receiverId,
      'sender_name': senderName,
      'receiver_name': receiverName,
    };
  }
}

class NotificationProperty {
  final int id;
  final String name;

  NotificationProperty({
    required this.id,
    required this.name,
  });

  factory NotificationProperty.fromJson(Map<String, dynamic> json) {
    return NotificationProperty(
      id: json['id'],
      name: json['name'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
    };
  }
}

class NotificationFloor {
  final int id;
  final String name;

  NotificationFloor({
    required this.id,
    required this.name,
  });

  factory NotificationFloor.fromJson(Map<String, dynamic> json) {
    return NotificationFloor(
      id: json['id'],
      name: json['name'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
    };
  }
} 