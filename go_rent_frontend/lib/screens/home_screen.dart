import 'package:flutter/material.dart';
import '../services/api_service.dart';
import '../services/notification_service.dart';
import '../models/property.dart';
import '../screens/property_details_screen.dart';
import '../screens/notifications_screen.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({Key? key}) : super(key: key);

  @override
  _HomeScreenState createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  final ApiService _apiService = ApiService();
  bool _isLoading = true;
  String? _error;
  List<Property> _properties = [];
  int _unreadNotifications = 0;

  @override
  void initState() {
    super.initState();
    _loadProperties();
    _loadNotifications();
  }

  Future<void> _loadNotifications() async {
    try {
      final notifications = await _apiService.getNotifications();
      setState(() {
        _unreadNotifications = notifications.where((n) => n.status == 'pending').length;
      });
    } catch (e) {
      print('Error loading notifications: $e');
    }
  }

  Future<void> _loadProperties() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      final properties = await _apiService.getProperties();
      setState(() {
        _properties = properties;
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
      if (e.toString().contains('Session expired')) {
        Navigator.of(context).pushReplacementNamed('/login');
      }
    }
  }

  Future<void> _testPushNotification() async {
    try {
      final response = await _apiService.testPushNotification(
        title: 'Test Push Notification',
        body: 'This is a test push notification from your app!',
        data: {'type': 'test', 'timestamp': DateTime.now().millisecondsSinceEpoch},
      );
      
      if (response) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Test push notification sent! Check your device.'),
            backgroundColor: Colors.green,
          ),
        );
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Error sending test notification: $e'),
          backgroundColor: Colors.red,
        ),
      );
    }
  }

  Future<void> _showAddPropertyDialog() async {
    final nameController = TextEditingController();
    final addressController = TextEditingController();

    return showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Add New Property'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: nameController,
              decoration: const InputDecoration(
                labelText: 'Property Name',
                hintText: 'Enter property name',
              ),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: addressController,
              decoration: const InputDecoration(
                labelText: 'Address',
                hintText: 'Enter property address',
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () async {
              if (nameController.text.isEmpty || addressController.text.isEmpty) {
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('Please fill all fields')),
                );
                return;
              }

              try {
                final success = await _apiService.addProperty(
                  nameController.text,
                  addressController.text,
                );

                if (success) {
                  Navigator.pop(context);
                  _loadProperties();
                } else {
                  throw Exception('Failed to add property');
                }
              } catch (e) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(content: Text('Error: $e')),
                );
              }
            },
            child: const Text('Add'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('My Properties'),
        actions: [
          // Test notification button
          IconButton(
            icon: const Icon(Icons.send, color: Colors.orange),
            onPressed: _testPushNotification,
            tooltip: 'Test Push Notification',
          ),
          Stack(
            children: [
              IconButton(
                icon: const Icon(Icons.notifications),
                onPressed: () async {
                      await Navigator.pushNamed(context, '/notifications');
                  _loadNotifications();
                },
              ),
              if (_unreadNotifications > 0)
                Positioned(
                  right: 8,
                  top: 8,
                  child: Container(
                    padding: const EdgeInsets.all(2),
                    decoration: BoxDecoration(
                      color: Colors.red,
                      borderRadius: BorderRadius.circular(10),
                    ),
                    constraints: const BoxConstraints(
                      minWidth: 16,
                      minHeight: 16,
                    ),
                    child: Text(
                      '$_unreadNotifications',
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 10,
                      ),
                      textAlign: TextAlign.center,
                    ),
                  ),
                ),
            ],
          ),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _error != null
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(_error!),
                      const SizedBox(height: 16),
                      ElevatedButton(
                        onPressed: _loadProperties,
                        child: const Text('Retry'),
                      ),
                    ],
                  ),
                )
              : Column(
                  children: [
                    // Always show the test notification button at the top
                    Container(
                      width: double.infinity,
                      padding: const EdgeInsets.all(16),
                      child: ElevatedButton(
                        onPressed: _testPushNotification,
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.orange,
                          foregroundColor: Colors.white,
                          padding: const EdgeInsets.symmetric(vertical: 16),
                          textStyle: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                        ),
                        child: const Text('🧪 TEST PUSH NOTIFICATION'),
                      ),
                    ),
                                        // Main content
                    Expanded(
                      child: _properties.isEmpty
                        ? Center(
                            child: Column(
                              mainAxisAlignment: MainAxisAlignment.center,
                              children: [
                                const Text('No properties added yet'),
                                const SizedBox(height: 16),
                                ElevatedButton(
                                  onPressed: _showAddPropertyDialog,
                                  child: const Text('Add First Property'),
                                ),
                              ],
                            ),
                          )
                        : ListView.builder(
                            padding: const EdgeInsets.symmetric(horizontal: 16),
                            itemCount: _properties.length,
                            itemBuilder: (context, index) {
                              final property = _properties[index];
                              return Card(
                                margin: const EdgeInsets.only(bottom: 16),
                                child: ListTile(
                                  leading: const Icon(Icons.home),
                                  title: Text(property.name),
                                  subtitle: Text(property.address),
                                  onTap: () {
                                    Navigator.push(
                                      context,
                                      MaterialPageRoute(
                                        builder: (context) => PropertyDetailsScreen(
                                          property: property,
                                        ),
                                      ),
                                    ).then((_) => _loadProperties());
                                  },
                                ),
                              );
                            },
                          ),
                    ),
                  ],
                ),
      floatingActionButton: FloatingActionButton(
        onPressed: _showAddPropertyDialog,
        child: const Icon(Icons.add),
      ),
    );
  }
} 