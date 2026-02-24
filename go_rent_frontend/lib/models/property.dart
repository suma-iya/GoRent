class Property {
  final int id;
  final String name;
  final String address;
  final String createdAt;
  final String? photo;

  Property({
    required this.id,
    required this.name,
    required this.address,
    required this.createdAt,
    this.photo,
  });

  factory Property.fromJson(Map<String, dynamic> json) {
    return Property(
      id: json['id'],
      name: json['name'],
      address: json['address'],
      createdAt: json['created_at'],
      photo: json['photo'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'address': address,
      'created_at': createdAt,
      'photo': photo,
    };
  }
} 