import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import '../models/property.dart';
import '../models/floor.dart';
import '../services/api_service.dart';
import 'pending_payment_notifications_screen.dart';
import 'payment_screen.dart';
import 'advance_details_screen.dart';
import '../utils/app_localizations.dart';

class PropertyDetailsScreen extends StatefulWidget {
  final Property property;
  final bool? isManagedProperty; // Add this parameter

  const PropertyDetailsScreen({Key? key, required this.property, this.isManagedProperty}) : super(key: key);

  @override
  _PropertyDetailsScreenState createState() => _PropertyDetailsScreenState();
}

class _PropertyDetailsScreenState extends State<PropertyDetailsScreen> {
  final ApiService _apiService = ApiService();
  bool _isLoading = true;
  String? _error;
  List<Floor> _floors = [];
  Property? _property;
  bool _isManager = false;

  // Modern color schemes
  final Color _managedColor = const Color(0xFF6366F1); // Modern indigo
  final Color _tenantColor = const Color(0xFF10B981); // Modern emerald
  final Color _backgroundColor = const Color(0xFFF8FAFC);
  final Color _cardColor = Colors.white;
  final Color _textPrimary = const Color(0xFF1E293B);
  final Color _textSecondary = const Color(0xFF64748B);
  
  Color get _themeColor => _isManager ? _managedColor : _tenantColor;
  String _getMonthName(BuildContext context, int month) {
  final local = AppLocalizations.of(context)!;
  final monthMap = {
  1: local.monthJanuary,
  2: local.monthFebruary,
  3: local.monthMarch,
  4: local.monthApril,
  5: local.monthMay,
  6: local.monthJune,
  7: local.monthJuly,
  8: local.monthAugust,
  9: local.monthSeptember,
  10: local.monthOctober,
  11: local.monthNovember,
  12: local.monthDecember,
  };
  return monthMap[month] ?? '';
}


  @override
  void initState() {
    super.initState();
    // Initialize with passed value to avoid color flash
    _isManager = widget.isManagedProperty ?? false;
    _apiService.loadCurrentUserId().then((_) {
      setState(() {}); // Rebuild to update UI with loaded userId
    });
    _loadPropertyDetails();
  }



  Future<void> _loadPropertyDetails() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      final result = await _apiService.getPropertyDetails(widget.property.id);
      
      setState(() {
        _property = result['property'];
        _floors = result['floors'];
        _isManager = result['is_manager'] ?? false;
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

  Future<void> _showAddFloorDialog() async {
    final nameController = TextEditingController();
    final rentController = TextEditingController();
    

    return showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(AppLocalizations.of(context).addFloor),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: nameController,
              decoration: InputDecoration(
                labelText: AppLocalizations.of(context).floorName,
                hintText: AppLocalizations.of(context).floorNameHint,
              ),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: rentController,
              decoration: InputDecoration(
                labelText: AppLocalizations.of(context).monthlyRent,
                hintText: AppLocalizations.of(context).rentHint,
              ),
              keyboardType: TextInputType.number,
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text(AppLocalizations.of(context).cancel),
          ),
          TextButton(
            onPressed: () async {
              if (nameController.text.isEmpty || rentController.text.isEmpty) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(content: Text(AppLocalizations.of(context).fillAllFields)),
                );
                return;
              }

              try {
                final rent = int.tryParse(rentController.text);
                if (rent == null) {
                  throw Exception(AppLocalizations.of(context).invalidRentAmount);
                }

                final success = await _apiService.addFloor(
                  widget.property.id,
                  nameController.text,
                  rent,
                );

                if (success) {
                  Navigator.pop(context);
                  _loadPropertyDetails();
                } else {
                  throw Exception(AppLocalizations.of(context).failedToAddFloor);
                }
              } catch (e) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(content: Text('${AppLocalizations.of(context).errorPrefix}: $e')),
                );
              }
            },
            child: Text(AppLocalizations.of(context).add),
          ),
        ],
      ),
    );
  }

  Future<void> _showAddTenantDialog(Floor floor) async {
    final nameController = TextEditingController();
    final phoneController = TextEditingController();

    return showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(AppLocalizations.of(context).addTenant),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: nameController,
                decoration: InputDecoration(
                  labelText: AppLocalizations.of(context).tenantName,
                  hintText: AppLocalizations.of(context).enterTenantName,
              ),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: phoneController,
                decoration: InputDecoration(
                  labelText: AppLocalizations.of(context).phoneNumber,
                  hintText: AppLocalizations.of(context).enterPhoneNumber,
              ),
              keyboardType: TextInputType.phone,
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text(AppLocalizations.of(context).cancel),
          ),
          TextButton(
            onPressed: () async {
              if (nameController.text.isEmpty || phoneController.text.isEmpty) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(content: Text(AppLocalizations.of(context).fillAllFields)),
                );
                return;
              }

              try {
                final success = await _apiService.addTenantToFloor(
                  widget.property.id,
                  floor.id,
                  nameController.text,
                  phoneController.text,
                );

                if (success) {
                  Navigator.pop(context);
                  _loadPropertyDetails();
                } else {
                  throw Exception(AppLocalizations.of(context).failedToAddTenant);
                }
              } catch (e) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(content: Text('${AppLocalizations.of(context).errorPrefix}: $e')),
                );
              }
            },
            child: Text(AppLocalizations.of(context).add),
          ),
        ],
      ),
    );
  }

  Future<void> _showRemoveTenantDialog(Floor floor) async {
    return showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(AppLocalizations.of(context).removeTenant),
        content: Text(AppLocalizations.of(context).removeTenantConfirmation),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text(AppLocalizations.of(context).cancel),
          ),
          TextButton(
            onPressed: () async {
              try {
                final success = await _apiService.removeTenantFromFloor(
                  widget.property.id,
                  floor.id,
                );

                if (success) {
                  Navigator.pop(context);
                  _loadPropertyDetails();
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text(AppLocalizations.of(context).tenantRemovedSuccessfully)),
                  );
                } else {
                  throw Exception(AppLocalizations.of(context).failedToRemoveTenant);
                }
              } catch (e) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(content: Text('${AppLocalizations.of(context).errorPrefix}: $e')),
                );
              }
            },
            child: Text(AppLocalizations.of(context).remove),
          ),
        ],
      ),
    );
  }

  Future<void> _showUpdateFloorDialog(Floor floor) async {
    final nameController = TextEditingController(text: floor.name);
    final rentController = TextEditingController(text: floor.rent.toString());

    return showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(AppLocalizations.of(context).updateFloor),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: nameController,
              decoration: InputDecoration(
                labelText: AppLocalizations.of(context).floorName,
                hintText: AppLocalizations.of(context).floorNameHint,
              ),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: rentController,
              decoration: InputDecoration(
                labelText: AppLocalizations.of(context).monthlyRent,
                hintText: AppLocalizations.of(context).rentHint,
              ),
              keyboardType: TextInputType.number,
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text(AppLocalizations.of(context).cancel),
          ),
          TextButton(
            onPressed: () async {
              if (nameController.text.isEmpty || rentController.text.isEmpty) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(content: Text(AppLocalizations.of(context).fillAllFields)),
                );
                return;
              }

              try {
                final rent = int.tryParse(rentController.text);
                if (rent == null) {
                  throw Exception(AppLocalizations.of(context).invalidRentAmount);
                }

                final success = await _apiService.updateFloor(
                  widget.property.id,
                  floor.id,
                  nameController.text,
                  rent,
                );

                if (success) {
                  Navigator.pop(context);
                  _loadPropertyDetails();
                } else {
                  throw Exception(AppLocalizations.of(context).failedToUpdateFloor);
                }
              } catch (e) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(content: Text('${AppLocalizations.of(context).errorPrefix}: $e')),
                );
              }
            },
            child: Text(AppLocalizations.of(context).update),
          ),
        ],
      ),
    );
  }

  void _showTenantRequestDialog(Floor floor) {
    final phoneController = TextEditingController();
    List<String> allPhones = [];
    List<String> filteredPhones = [];
    bool isLoading = false;

    showDialog(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setState) {
          return AlertDialog(
            title: Text(AppLocalizations.of(context).sendTenantRequest),
            content: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                if (isLoading)
                  const Center(child: CircularProgressIndicator())
                else
                  TextField(
                    controller: phoneController,
                    decoration: InputDecoration(
                      labelText: AppLocalizations.of(context).phoneNumber,
                      suffixIcon: IconButton(
                        icon: const Icon(Icons.clear),
                        onPressed: () {
                          phoneController.clear();
                          setState(() {
                            filteredPhones = [];
                          });
                        },
                      ),
                    ),
                    onTap: () async {
                      if (allPhones.isEmpty) {
                        setState(() {
                          isLoading = true;
                        });
                        try {
                          allPhones = await _apiService.getUserPhones();
                          setState(() {
                            filteredPhones = allPhones;
                            isLoading = false;
                          });
                        } catch (e) {
                          setState(() {
                            isLoading = false;
                          });
                          if (mounted) {
                            ScaffoldMessenger.of(context).showSnackBar(
                              SnackBar(content: Text('${AppLocalizations.of(context).errorPrefix}: $e')),
                            );
                          }
                        }
                      }
                    },
                    onChanged: (value) {
                      setState(() {
                        filteredPhones = allPhones
                            .where((phone) => phone.contains(value))
                            .toList();
                      });
                    },
                  ),
                if (filteredPhones.isNotEmpty)
                  Container(
                    constraints: const BoxConstraints(maxHeight: 200),
                    child: ListView.builder(
                      shrinkWrap: true,
                      itemCount: filteredPhones.length,
                      itemBuilder: (context, index) {
                        return ListTile(
                          title: Text(filteredPhones[index]),
                          onTap: () {
                            phoneController.text = filteredPhones[index];
                            setState(() {
                              filteredPhones = [];
                            });
                          },
                        );
                      },
                    ),
                  ),
              ],
            ),
            actions: [
              TextButton(
                onPressed: () => Navigator.pop(context),
                child: Text(AppLocalizations.of(context).cancel),
              ),
              TextButton(
                onPressed: () async {
                  if (phoneController.text.isEmpty) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(content: Text(AppLocalizations.of(context).pleaseEnterPhoneNumber)),
                    );
                    return;
                  }

                  try {
                    final success = await _apiService.sendTenantRequest(
                      widget.property.id,
                      floor.id,
                      phoneController.text,
                    );

                    if (success) {
                      if (mounted) {
                        Navigator.pop(context);
                        ScaffoldMessenger.of(context).showSnackBar(
                          SnackBar(content: Text(AppLocalizations.of(context).tenantRequestSentSuccessfully)),
                        );
                        await _loadPropertyDetails(); // Refresh property details
                      }
                    }
                  } catch (e) {
                    if (mounted) {
                      // Show error message first
                      ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(
                          content: Text(e.toString().replaceAll('Exception: ', '')),
                          backgroundColor: Colors.red,
                          duration: const Duration(seconds: 3),
                        ),
                      );
                      
                      // Wait for the snackbar to be visible
                      await Future.delayed(const Duration(milliseconds: 500));
                      
                      // Then close the dialog and refresh
                      Navigator.pop(context);
                      await _loadPropertyDetails();
                    }
                  }
                },
                child: Text(AppLocalizations.of(context).sendRequest),
              ),
            ],
          );
        },
      ),
    );
  }

  Future<void> _showCancelRequestDialog(Floor floor) async {
    return showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(AppLocalizations.of(context).cancelRequest),
        content: Text(AppLocalizations.of(context).cancelRequestConfirmation),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text(AppLocalizations.of(context).no),
          ),
          TextButton(
            onPressed: () async {
              try {
                final success = await _apiService.deleteNotification(floor.notificationId!);
                if (success) {
                  Navigator.pop(context);
                  _loadPropertyDetails();
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text(AppLocalizations.of(context).requestCancelledSuccessfully)),
                  );
                } else {
                  throw Exception(AppLocalizations.of(context).failedToCancelRequest);
                }
              } catch (e) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(content: Text('${AppLocalizations.of(context).errorPrefix}: $e')),
                );
              }
            },
            child: Text(AppLocalizations.of(context).yes),
          ),
        ],
      ),
    );
  }

  Future<void> _showSendPaymentDialog(Floor floor) async {
    final amountController = TextEditingController();
    final electricityBillController = TextEditingController();
    
    return showDialog(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setState) => AlertDialog(
        title: Text(AppLocalizations.of(context).sendPaymentNotification),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
          controller: amountController,
          decoration: InputDecoration(
            labelText: AppLocalizations.of(context).paymentAmount,
            hintText: AppLocalizations.of(context).enterPaymentAmount,
          ),
          keyboardType: TextInputType.number,
              ),
              const SizedBox(height: 16),
              TextField(
                controller: electricityBillController,
                                  decoration: InputDecoration(
                    labelText: AppLocalizations.of(context).electricityBill,
                    hintText: AppLocalizations.of(context).enterElectricityBill,
                  ),
                keyboardType: TextInputType.number,
              ),
            ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text(AppLocalizations.of(context).cancel),
          ),
          TextButton(
            onPressed: () async {
              final amount = int.tryParse(amountController.text);
              if (amount == null) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(content: Text(AppLocalizations.of(context).pleaseEnterValidAmount)),
                );
                return;
              }
              
              // Parse electricity bill if provided
              int? electricityBill;
              if (electricityBillController.text.isNotEmpty) {
                electricityBill = int.tryParse(electricityBillController.text);
                if (electricityBill == null) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text(AppLocalizations.of(context).pleaseEnterValidElectricityBill)),
                  );
                  return;
                }
              }
              
              try {
                await _apiService.sendPaymentNotification(
                  propertyId: widget.property.id,
                  floorId: floor.id,
                  amount: amount,
                  electricityBill: electricityBill,
                );
                Navigator.pop(context);
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(content: Text(AppLocalizations.of(context).paymentNotificationSent)),
                );
              } catch (e) {
                Navigator.pop(context);
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(content: Text('Error: $e')),
                );
              }
            },
            child: Text(AppLocalizations.of(context).send),
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _showAdvancePaymentDialog(Floor floor) async {
    final phoneController = TextEditingController();
    final moneyController = TextEditingController();
    List<String> allPhones = [];
    List<String> filteredPhones = [];
    bool isLoading = false;
    
    return showDialog(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setState) {
          return AlertDialog(
              title: Text(AppLocalizations.of(context).sendAdvancePaymentRequest),
              content: SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  // Phone number field with search
                  if (isLoading)
                    const Center(child: CircularProgressIndicator())
                  else
                    TextField(
                      controller: phoneController,
                      decoration: InputDecoration(
                        labelText: AppLocalizations.of(context).phoneNumber,
                        hintText: AppLocalizations.of(context).enterPhoneNumber,
                        suffixIcon: IconButton(
                          icon: const Icon(Icons.clear),
                          onPressed: () {
                            phoneController.clear();
                            setState(() {
                              filteredPhones = [];
                            });
                          },
                        ),
                      ),
                      onTap: () async {
                        if (allPhones.isEmpty) {
                          setState(() {
                            isLoading = true;
                          });
                          try {
                            allPhones = await _apiService.getUserPhones();
                            setState(() {
                              filteredPhones = allPhones;
                              isLoading = false;
                            });
                          } catch (e) {
                            setState(() {
                              isLoading = false;
                            });
                            if (mounted) {
                              ScaffoldMessenger.of(context).showSnackBar(
                                SnackBar(content: Text('${AppLocalizations.of(context).errorPrefix}: $e')),
                              );
                            }
                          }
                        }
                      },
                      onChanged: (value) {
                        setState(() {
                          filteredPhones = allPhones
                              .where((phone) => phone.contains(value))
                              .toList();
                        });
                      },
                    ),
                  if (filteredPhones.isNotEmpty)
                    Container(
                      constraints: const BoxConstraints(maxHeight: 120),
                      child: ListView.builder(
                        shrinkWrap: true,
                        itemCount: filteredPhones.length,
                        itemBuilder: (context, index) {
                          return ListTile(
                            dense: true,
                            title: Text(
                              filteredPhones[index],
                              style: const TextStyle(fontSize: 14),
                            ),
                            onTap: () {
                              phoneController.text = filteredPhones[index];
                              setState(() {
                                filteredPhones = [];
                              });
                            },
                          );
                        },
                      ),
                    ),
                  const SizedBox(height: 12),
                  // Amount field
                  TextField(
                    controller: moneyController,
                    decoration: InputDecoration(
                      labelText: AppLocalizations.of(context).amount,
                      hintText: AppLocalizations.of(context).enterAdvancePaymentAmount,
                    ),
                    keyboardType: TextInputType.number,
                  ),

                ],
              ),
            ),
            actions: [
              TextButton(
                onPressed: () => Navigator.pop(context),
                child: Text(AppLocalizations.of(context).cancel),
              ),
              TextButton(
                onPressed: () async {
                  if (phoneController.text.isEmpty) {
                    ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text(AppLocalizations.of(context).pleaseEnterPhoneNumber)),
                    );
                    return;
                  }
                  
                  final money = int.tryParse(moneyController.text);
                  if (money == null || money <= 0) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(content: Text(AppLocalizations.of(context).pleaseEnterValidAmount)),
                    );
                    return;
                  }
                  

                  
                  try {
                    // Get user ID from phone number
                    final userId = await _apiService.getUserIdByPhone(phoneController.text);
                    
                    final success = await _apiService.createAdvancePaymentRequest(
                      propertyId: widget.property.id,
                      floorId: floor.id,
                      advanceUid: userId,
                      money: money,
                    );
                    
                    if (success) {
                      Navigator.pop(context);
                      ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(content: Text(AppLocalizations.of(context).advancePaymentRequestSentSuccessfully)),
                      );
                      // Reload property details to update the UI
                      _loadPropertyDetails();
                    } else {
                      throw Exception('Failed to send advance payment request');
                    }
                  } catch (e) {
                    Navigator.pop(context);
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(content: Text('${AppLocalizations.of(context).errorPrefix}: $e')),
                    );
                  }
                },
                child: Text(AppLocalizations.of(context).send),
              ),
            ],
          );
        },
      ),
    );
  }

  Future<void> _showCancelAdvancePaymentDialog(Floor floor) async {
    return showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(AppLocalizations.of(context).cancelAdvancePaymentRequest),
        content: Text(AppLocalizations.of(context).cancelAdvancePaymentRequestConfirmation),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text(AppLocalizations.of(context).no),
          ),
          TextButton(
            onPressed: () async {
              try {
                final success = await _apiService.cancelAdvancePaymentRequest(floor.id);
                
                if (success) {
                  Navigator.pop(context);
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text(AppLocalizations.of(context).advancePaymentRequestCancelledSuccessfully)),
                  );
                  // Reload property details to update the UI
                  _loadPropertyDetails();
                } else {
                  throw Exception(AppLocalizations.of(context).failedToCancelAdvancePaymentRequest);
                }
              } catch (e) {
                Navigator.pop(context);
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(content: Text('${AppLocalizations.of(context).errorPrefix}: $e')),
                );
              }
            },
            child: Text(AppLocalizations.of(context).yesCancelAdvancePaymentRequest),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: _backgroundColor,
      appBar: AppBar(
        title: Text(
          _property?.name ?? widget.property.name,
          style: const TextStyle(
            fontWeight: FontWeight.w700,
            fontSize: 20,
            color: Colors.white,
          ),
        ),
        centerTitle: true,
        elevation: 0,
        backgroundColor: _themeColor,
        systemOverlayStyle: SystemUiOverlayStyle.light,
        leading: Hero(
          tag: 'property_${widget.property.id}',
          child: Container(
            margin: const EdgeInsets.all(8),
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.2),
              shape: BoxShape.circle,
            ),
            child: IconButton(
              icon: Icon(
                _isManager ? Icons.business_rounded : Icons.person_rounded,
                color: Colors.white,
              ),
              onPressed: () => Navigator.pop(context),
            ),
          ),
        ),
        actions: [],
      ),
      body: SafeArea(
        child: _isLoading
            ? const Center(child: CircularProgressIndicator())
            : _error != null
                ? Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Text(_error!),
                        const SizedBox(height: 16),
                        ElevatedButton(
                          onPressed: _loadPropertyDetails,
                          child: Text(AppLocalizations.of(context).retry),
                          style: ElevatedButton.styleFrom(
                            backgroundColor: _themeColor,
                            foregroundColor: Colors.white,
                            padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 14),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(12),
                            ),
                            elevation: 0,
                          ),
                        ),
                      ],
                    ),
                  )
                : _floors.isEmpty
                    ? Center(
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Text(AppLocalizations.of(context).noFloorsAddedYet),
                            if (_isManager) ...[
                              const SizedBox(height: 16),
                              ElevatedButton(
                                onPressed: _showAddFloorDialog,
                                child: Text(AppLocalizations.of(context).addFirstFloor),
                                style: ElevatedButton.styleFrom(
                                  backgroundColor: _themeColor,
                                  foregroundColor: Colors.white,
                                  padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 14),
                                  shape: RoundedRectangleBorder(
                                    borderRadius: BorderRadius.circular(12),
                                  ),
                                  elevation: 0,
                                ),
                              ),
                            ],
                          ],
                        ),
                      )
                    : ListView.builder(
                        padding: const EdgeInsets.all(20),
                        itemCount: _floors.length,
                        itemBuilder: (context, index) {
                          final floor = _floors[index];
                          // Determine status: if tenant exists and no explicit status, it's occupied
                          final actualStatus = floor.status ?? (floor.tenant != null ? 'occupied' : 'available');
                          final themeColor = actualStatus == 'occupied'
                              ? _themeColor
                              : actualStatus == 'pending'
                                  ? Colors.orange
                                  : Color.fromARGB(255, 118, 221, 235); // Make available status green
                          return Container(
                            margin: const EdgeInsets.only(bottom: 20),
                            decoration: BoxDecoration(
                              color: _cardColor,
                              borderRadius: BorderRadius.circular(20),
                              boxShadow: [
                                BoxShadow(
                                  color: Colors.black.withOpacity(0.08),
                                  blurRadius: 20,
                                  spreadRadius: 0,
                                  offset: const Offset(0, 8),
                                ),
                              ],
                              border: Border.all(
                                color: themeColor.withOpacity(0.1),
                                width: 1,
                              ),
                            ),
                            child: Material(
                              color: Colors.transparent,
                              child: InkWell(
                                borderRadius: BorderRadius.circular(20),
                                onTap: () {
                                  // Navigate to payment screen when floor is clicked
                                  Navigator.push(
                                    context,
                                    MaterialPageRoute(
                                      builder: (context) => PaymentScreen(
                                        property: widget.property,
                                        floor: floor,
                                        themeColor: _themeColor,
                                      ),
                                    ),
                                  );
                                },
                            child: Column(
                              children: [
                                    Padding(
                                      padding: const EdgeInsets.all(20),
                                      child: Row(
                                        crossAxisAlignment: CrossAxisAlignment.start,
                                        children: [
                                          Container(
                                            padding: const EdgeInsets.all(16),
                                            decoration: BoxDecoration(
                                              gradient: LinearGradient(
                                                colors: [
                                                  themeColor.withOpacity(0.2),
                                                  themeColor.withOpacity(0.1),
                                                ],
                                                begin: Alignment.topLeft,
                                                end: Alignment.bottomRight,
                                              ),
                                              borderRadius: BorderRadius.circular(16),
                                            ),
                                            child: Icon(
                                              actualStatus == 'occupied'
                                                  ? Icons.person_rounded
                                                  : actualStatus == 'pending'
                                                      ? Icons.hourglass_top_rounded
                                                      : Icons.home_rounded,
                                              color: themeColor,
                                              size: 28,
                                  ),
                                          ),
                                          const SizedBox(width: 16),
                                          Expanded(
                                            child: Column(
                                    crossAxisAlignment: CrossAxisAlignment.start,
                                    children: [
                                        Row(
                                          children: [
                                                    Expanded(
                                                      child: Text(
                                                        floor.name,
                                                        style: TextStyle(
                                                          fontSize: 18,
                                                          fontWeight: FontWeight.w700,
                                                          color: _textPrimary,
                                                        ),
                                                        overflow: TextOverflow.ellipsis,
                                                      ),
                                                    ),
                                              const SizedBox(width: 8),
                                                    Container(
                                                      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                                                      decoration: BoxDecoration(
                                                        color: themeColor.withOpacity(0.12),
                                                        borderRadius: BorderRadius.circular(12),
                                                      ),
                                                                                                                child: Text(
                                                            actualStatus == 'occupied'
                                                               ? AppLocalizations.of(context).occupied
                                                               : actualStatus == 'pending'
                                                                 ? AppLocalizations.of(context).pending
                                                                 : AppLocalizations.of(context).available,
                                                        style: TextStyle(
                                                          fontSize: 12,
                                                          fontWeight: FontWeight.w600,
                                                          color: themeColor,
                                                        ),
                                                      ),
                                                    ),
                                                  ],
                                                ),
                                                const SizedBox(height: 6),
                                                // Rent information
                                                Row(
                                                  children: [
                                                    Icon(
                                                      Icons.attach_money_rounded,
                                                      size: 16,
                                                      color: _textSecondary,
                                                    ),
                                                    const SizedBox(width: 4),
                                                    Expanded(
                                                      child: Text(
                                                        '${AppLocalizations.of(context).rent}: ${floor.rent} ${AppLocalizations.of(context).currency}',
                                                        style: TextStyle(
                                                          fontSize: 14,
                                                          color: _textSecondary,
                                                        ),
                                                        overflow: TextOverflow.ellipsis,
                                                      ),
                                                    ),
                                                  ],
                                                ),
                                                // Tenant information
                                                if (floor.tenant != null) ...[
                                                  const SizedBox(height: 4),
                                                  Row(
                                                    children: [
                                                      Icon(
                                                        Icons.person,
                                                        size: 16,
                                                        color: _textSecondary,
                                                      ),
                                                      const SizedBox(width: 4),
                                                      Expanded(
                                                        child: Text(
                                                          '${AppLocalizations.of(context).tenant}: ${floor.tenantName ?? 'Unknown User'}',
                                                          style: TextStyle(
                                                            fontSize: 14,
                                                            color: _textSecondary,
                                                          ),
                                                          overflow: TextOverflow.ellipsis,
                                                        ),
                                                      ),
                                                    ],
                                                  ),
                                                ],
                                                // Advance details button
                                                const SizedBox(height: 8),
                                                GestureDetector(
                                                  onTap: () {
                                                    Navigator.push(
                                                      context,
                                                      MaterialPageRoute(
                                                        builder: (context) => AdvanceDetailsScreen(
                                                          floorName: floor.name,
                                                          floorId: floor.id,
                                                          themeColor: themeColor,
                                                        ),
                                                      ),
                                                    );
                                                  },
                                                  child: Row(
                                                    children: [
                                                      Icon(
                                                        Icons.account_balance_wallet_rounded,
                                                        size: 16,
                                                        color: themeColor,
                                                      ),
                                                      const SizedBox(width: 4),
                                                      Expanded(
                                                        child: Text(
                                                          AppLocalizations.of(context).showAdvanceDetails,
                                                          style: TextStyle(
                                                            fontSize: 14,
                                                            color: themeColor,
    
                                                          ),
                                                          overflow: TextOverflow.ellipsis,
                                                        ),
                                                      ),
                                                    ],
                                                  ),
                                                ),
                                                if (actualStatus == 'pending') ...[
                                                  const SizedBox(height: 4),
                                                  Row(
                                                    children: [
                                                      const Icon(Icons.hourglass_top_rounded, size: 16, color: Colors.orange),
                                                      const SizedBox(width: 4),
                                                      Expanded(
                                                        child: Text(
                                                          AppLocalizations.of(context).requestPending,
                                                          style: TextStyle(
                                                            color: Colors.orange,
                                                            fontWeight: FontWeight.bold,
                                                            fontSize: 13,
                                                          ),
                                                          overflow: TextOverflow.ellipsis,
                                                        ),
                                                      ),
                                                      if (floor.tenant != null && floor.tenant == _apiService.currentUserId) ...[
                                                        const SizedBox(width: 8),
                                                        const Icon(Icons.touch_app, size: 16, color: Colors.orange),
                                                        GestureDetector(
                                  onTap: () {
                                      Navigator.push(
                                        context,
                                        MaterialPageRoute(
                                          builder: (context) => PendingPaymentNotificationsScreen(
                                            property: widget.property,
                                            floor: floor,
                                          ),
                                        ),
                                      );
                                                          },
                                                          child: Text(AppLocalizations.of(context).tapToView,
                                                            style: TextStyle(color: Colors.orange, fontSize: 12)),
                                                        ),
                                                      ],
                                                    ],
                                                  ),
                                                ],
                                              ],
                                            ),
                                          ),
                                          const SizedBox(width: 8),
                                          if (_isManager)
                                            Column(
                                    children: [
                                      IconButton(
                                        icon: const Icon(Icons.edit, color: Colors.blue),
                                        onPressed: () => _showUpdateFloorDialog(floor),
                                      ),
                                      if (floor.tenant != null)
                                        IconButton(
                                          icon: const Icon(Icons.person_remove, color: Colors.red),
                                          onPressed: () => _showRemoveTenantDialog(floor),
                                        )
                                                else if (actualStatus == 'pending' && floor.notificationId != null)
                                        IconButton(
                                          icon: const Icon(Icons.cancel, color: Colors.orange),
                                          onPressed: () => _showCancelRequestDialog(floor),
                                        )
                                      else
                                        IconButton(
                                          icon: const Icon(Icons.person_add, color: Colors.green),
                                          onPressed: () => _showTenantRequestDialog(floor),
                                        ),
                                      // Advance payment icon for managers
                                      if (floor.hasPendingAdvancePayment)
                                        IconButton(
                                          icon: const Icon(Icons.cancel, color: Colors.red),
                                          onPressed: () => _showCancelAdvancePaymentDialog(floor),
                                        )
                                      else
                                        IconButton(
                                          icon: const Icon(Icons.account_balance_wallet_rounded, color: Colors.orange),
                                          onPressed: () => _showAdvancePaymentDialog(floor),
                                        ),
                                    ],
                                ),
                                        ],
                                      ),
                                ),
                                    // Payment button for tenants
                                if (floor.tenant != null && floor.tenant == _apiService.currentUserId)
                                  Padding(
                                        padding: const EdgeInsets.only(left: 20, right: 20, bottom: 16),
                                    child: SizedBox(
                                      width: double.infinity,
                                      child: ElevatedButton.icon(
                                            icon: const Icon(Icons.payment_rounded),
                                        label: Text(AppLocalizations.of(context).sendPaymentNotification),
                                            style: ElevatedButton.styleFrom(
                                              backgroundColor: _themeColor,
                                              foregroundColor: Colors.white,
                                              padding: const EdgeInsets.symmetric(vertical: 16),
                                              shape: RoundedRectangleBorder(
                                                borderRadius: BorderRadius.circular(12),
                                              ),
                                              elevation: 0,
                                            ),
                                        onPressed: () {
                                          _showSendPaymentDialog(floor);
                                        },
                                      ),
                                    ),
                                  ),
                              ],
                                ),
                              ),
                            ),
                          );
                        },
                      ),
      ),
      floatingActionButton: _isManager ? FloatingActionButton.extended(
        onPressed: _showAddFloorDialog,
        backgroundColor: _themeColor,
        foregroundColor: Colors.white,
        elevation: 8,
        icon: const Icon(Icons.add_rounded),
        label: Text(
          AppLocalizations.of(context).addFloor,
          style: TextStyle(fontWeight: FontWeight.w600),
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(16),
        ),
      ) : null,
    );
  }
} 