import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:intl/intl.dart';
import '../models/property.dart';
import '../models/floor.dart';
import '../services/api_service.dart';
import '../utils/app_localizations.dart';

class PaymentScreen extends StatefulWidget {
  final Property property;
  final Floor floor;
  final Color themeColor;

  const PaymentScreen({
    Key? key,
    required this.property,
    required this.floor,
    required this.themeColor,
  }) : super(key: key);

  @override
  _PaymentScreenState createState() => _PaymentScreenState();
}

class _PaymentScreenState extends State<PaymentScreen> {
  final ApiService _apiService = ApiService();
  bool _isLoading = true;
  String? _error;
  int _dueRent = 0;
  bool _isManager = false;
  List<Map<String, dynamic>> _payments = [];
  bool _isLoadingPayments = false;
  
  // Pagination variables
  int _currentPage = 1;
  int _totalPages = 1;
  int _totalCount = 0;
  bool _hasNextPage = false;
  bool _hasPrevPage = false;

  // Modern color schemes
  Color get _primaryColor => widget.themeColor;
  final Color _backgroundColor = const Color(0xFFF8FAFC);
  final Color _cardColor = Colors.white;
  final Color _textPrimary = const Color(0xFF1E293B);
  final Color _textSecondary = const Color(0xFF64748B);

  String _formatDate(String dateStr) {
    try {
      final dt = DateTime.parse(dateStr);
      return DateFormat('dd MMM yyyy, hh:mm a').format(dt);
    } catch (_) {
      return dateStr;
    }
  }

  Future<void> _loadPage(int page) async {
    setState(() {
      _isLoadingPayments = true;
    });

    try {
      final paymentHistoryData = await _apiService.getPaymentHistory(widget.floor.id, page: page);
      
      setState(() {
        _payments = (paymentHistoryData['payments'] as List<dynamic>?)?.cast<Map<String, dynamic>>() ?? [];
        
        // Update pagination info
        final pagination = paymentHistoryData['pagination'] as Map<String, dynamic>? ?? {};
        _currentPage = pagination['current_page'] ?? 1;
        _totalPages = pagination['total_pages'] ?? 1;
        _totalCount = pagination['total_count'] ?? 0;
        _hasNextPage = pagination['has_next_page'] ?? false;
        _hasPrevPage = pagination['has_prev_page'] ?? false;
        
        _isLoadingPayments = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoadingPayments = false;
      });
    }
  }

  @override
  void initState() {
    super.initState();
    _loadPaymentDetails();
  }

  Future<void> _loadPaymentDetails() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      final paymentDetails = await _apiService.getPaymentDetails(widget.floor.id);
      final isManager = await _apiService.checkUserManager(widget.property.id);
      final paymentHistoryData = await _apiService.getPaymentHistory(widget.floor.id, page: _currentPage);
      
      setState(() {
        _dueRent = paymentDetails['rent'] ?? 0;
        _isManager = isManager;
        _payments = (paymentHistoryData['payments'] as List<dynamic>?)?.cast<Map<String, dynamic>>() ?? [];
        
        // Update pagination info
        final pagination = paymentHistoryData['pagination'] as Map<String, dynamic>? ?? {};
        _currentPage = pagination['current_page'] ?? 1;
        _totalPages = pagination['total_pages'] ?? 1;
        _totalCount = pagination['total_count'] ?? 0;
        _hasNextPage = pagination['has_next_page'] ?? false;
        _hasPrevPage = pagination['has_prev_page'] ?? false;
        
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
    }
  }



  Future<void> _showCreatePaymentDialog() async {
    final amountController = TextEditingController();
    final electricityBillController = TextEditingController();
    bool isAddition = true; // true for add, false for subtract

    return showDialog(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setState) => AlertDialog(
          title: Text(AppLocalizations.of(context).adjustDueRent),
          content: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                // Operation selection (Add/Subtract)
                Row(
                  children: [
                    Expanded(
                      child: GestureDetector(
                        onTap: () => setState(() => isAddition = true),
                        child: Container(
                          padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 8),
                          decoration: BoxDecoration(
                            color: isAddition ? _primaryColor : Colors.grey.shade200,
                            borderRadius: BorderRadius.circular(8),
                          ),
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Icon(
                                Icons.add_rounded,
                                color: isAddition ? Colors.white : Colors.grey.shade600,
                                size: 18,
                              ),
                              const SizedBox(width: 4),
                              Flexible(
                                child: Text(
                                  AppLocalizations.of(context).add,
                                  style: TextStyle(
                                    color: isAddition ? Colors.white : Colors.grey.shade600,
                                    fontWeight: FontWeight.w600,
                                    fontSize: 14,
                                  ),
                                  overflow: TextOverflow.ellipsis,
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),
                    const SizedBox(width: 8),
                    Expanded(
                      child: GestureDetector(
                        onTap: () => setState(() => isAddition = false),
                        child: Container(
                          padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 8),
                          decoration: BoxDecoration(
                            color: !isAddition ? Colors.red.shade500 : Colors.grey.shade200,
                            borderRadius: BorderRadius.circular(8),
                          ),
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Icon(
                                Icons.remove_rounded,
                                color: !isAddition ? Colors.white : Colors.grey.shade600,
                                size: 18,
                              ),
                              const SizedBox(width: 4),
                              Flexible(
                                child: Text(
                                  AppLocalizations.of(context).subtract,
                                  style: TextStyle(
                                    color: !isAddition ? Colors.white : Colors.grey.shade600,
                                    fontWeight: FontWeight.w600,
                                    fontSize: 14,
                                  ),
                                  overflow: TextOverflow.ellipsis,
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              const SizedBox(height: 20),
              // Amount input
              TextField(
                controller: amountController,
                onChanged: (value) {
                  setState(() {}); // Trigger rebuild to update preview
                },
                decoration: InputDecoration(
                                      labelText: AppLocalizations.of(context).rentAmount,
                  hintText: '${AppLocalizations.of(context).enterPaymentAmount} ${isAddition ? AppLocalizations.of(context).add.toLowerCase() : AppLocalizations.of(context).subtract.toLowerCase()}',
                  prefixIcon: Icon(
                    isAddition ? Icons.add_rounded : Icons.remove_rounded,
                    color: isAddition ? _primaryColor : Colors.red.shade500,
                  ),
                ),
                keyboardType: TextInputType.number,
              ),
              const SizedBox(height: 16),
              // Electricity bill input
              TextField(
                controller: electricityBillController,
                decoration: InputDecoration(
                  labelText: AppLocalizations.of(context).electricityBillOptional,
                  hintText: AppLocalizations.of(context).enterElectricityBill,
                  prefixIcon: Icon(
                    Icons.electric_bolt_rounded,
                    color: Colors.orange.shade600,
                  ),
                ),
                keyboardType: TextInputType.number,
              ),
              const SizedBox(height: 16),
              // Preview of new due rent
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: Colors.grey.shade50,
                  borderRadius: BorderRadius.circular(8),
                  border: Border.all(color: Colors.grey.shade300),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Icon(
                          Icons.info_outline_rounded,
                          color: _textSecondary,
                          size: 20,
                        ),
                        const SizedBox(width: 8),
                        Text(
                          'Preview:',
                          style: TextStyle(
                            color: _textSecondary,
                            fontSize: 14,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 8),
                                          Text(
                        '${AppLocalizations.of(context).current}: ${_dueRent.toStringAsFixed(2)} ${AppLocalizations.of(context).tk}',
                      style: TextStyle(
                        color: _textSecondary,
                        fontSize: 14,
                      ),
                    ),
                    Text(
                      '${isAddition ? '+' : '-'} \$${(int.tryParse(amountController.text) ?? 0).toStringAsFixed(2)}',
                      style: TextStyle(
                        color: isAddition ? _primaryColor : Colors.red.shade500,
                        fontSize: 14,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    const SizedBox(height: 4),
                                          Text(
                        '${AppLocalizations.of(context).newTotal}: ${(isAddition ? _dueRent + (int.tryParse(amountController.text) ?? 0) : _dueRent - (int.tryParse(amountController.text) ?? 0)).toStringAsFixed(2)} ${AppLocalizations.of(context).tk}',
                      style: TextStyle(
                        color: _textPrimary,
                        fontSize: 16,
                        fontWeight: FontWeight.w700,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: Text(AppLocalizations.of(context).cancel),
            ),
            ElevatedButton(
              onPressed: () async {
                final amount = int.tryParse(amountController.text);
                
                if (amount == null || amount <= 0) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text(AppLocalizations.of(context).pleaseEnterValidPositiveAmount)),
                  );
                  return;
                }

                // Calculate new due rent
                final newDueRent = isAddition ? _dueRent + amount : _dueRent - amount;
                
                // Prevent negative due rent
                if (newDueRent < 0) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text(AppLocalizations.of(context).dueRentCannotBeNegative)),
                  );
                  return;
                }

                try {
                  // Parse electricity bill values
                  int? electricityBill;
                  
                  if (electricityBillController.text.isNotEmpty) {
                    electricityBill = int.tryParse(electricityBillController.text);
                  }
                  
                  final success = await _apiService.createPayment(
                    propertyId: widget.property.id,
                    floorId: widget.floor.id,
                    dueRent: amount, // Send only the amount being added/subtracted
                    receivedMoney: 0, // No received money for adjustments
                    fullPayment: false,
                    electricityBill: electricityBill,
                  );

                  if (success) {
                    Navigator.pop(context);
                    _loadPaymentDetails();
                    ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(content: Text('${AppLocalizations.of(context).dueRent} ${isAddition ? AppLocalizations.of(context).increased : AppLocalizations.of(context).decreased} ${AppLocalizations.of(context).successfully}!')),
                      );
                  } else {
                    throw Exception('Failed to update due rent');
                  }
                } catch (e) {
                  Navigator.pop(context);
                  // Extract clean error message
                  String errorMessage = e.toString();
                  if (errorMessage.startsWith('Exception: ')) {
                    errorMessage = errorMessage.substring(11);
                  }
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: Text(errorMessage),
                      backgroundColor: Colors.red,
                      duration: const Duration(seconds: 4),
                    ),
                  );
                }
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: isAddition ? _primaryColor : Colors.red.shade500,
                foregroundColor: Colors.white,
              ),
              child: Text(isAddition ? AppLocalizations.of(context).addAmount : AppLocalizations.of(context).subtractAmount),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildPaymentHistoryItem(Map<String, dynamic> payment) {
    final dueRent = (payment['rent'] as num?)?.toDouble() ?? 0.0;
    final month = payment['month'] as String? ?? '';
    final receivedMoney = (payment['received_money'] as num?)?.toDouble() ?? 0.0;
    final createdAt = payment['created_at'] as String? ?? '';
    final fullPayment = payment['full_payment'] as bool? ?? false;

    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      decoration: BoxDecoration(
        color: _cardColor,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            spreadRadius: 0,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Expanded(
                  child: Row(
                    children: [
                      Container(
                        padding: const EdgeInsets.all(8),
                        decoration: BoxDecoration(
                          color: fullPayment ? Colors.green.withOpacity(0.1) : Colors.orange.withOpacity(0.1),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Icon(
                          fullPayment ? Icons.check_circle_rounded : Icons.pending_rounded,
                          color: fullPayment ? Colors.green.shade600 : Colors.orange.shade600,
                          size: 20,
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              fullPayment ? 'Payment Completed' : 'Payment Pending',
                              style: TextStyle(
                                fontSize: 14,
                                fontWeight: FontWeight.w600,
                                color: fullPayment ? Colors.green.shade600 : Colors.orange.shade600,
                              ),
                              overflow: TextOverflow.ellipsis,
                            ),
                            if (month.isNotEmpty)
                              Text(
                                month,
                                style: TextStyle(
                                  fontSize: 12,
                                  color: _textSecondary,
                                ),
                                overflow: TextOverflow.ellipsis,
                              ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(width: 8),
                Flexible(
                  child: Text(
                    createdAt,
                    style: TextStyle(
                      fontSize: 12,
                      color: _textSecondary,
                    ),
                    textAlign: TextAlign.end,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: _buildPaymentDetail(
                                          AppLocalizations.of(context).dueRent,
                    '${dueRent.toStringAsFixed(2)} ${AppLocalizations.of(context).tk}',
                    Colors.red.shade600,
                    Icons.money_off_rounded,
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: _buildPaymentDetail(
                                          AppLocalizations.of(context).receivedMoney,
                    '${receivedMoney.toStringAsFixed(2)} ${AppLocalizations.of(context).tk}',
                    Colors.green.shade600,
                    Icons.money_rounded,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildPaymentHistoryTable() {
    if (_isLoadingPayments) {
      return Center(
        child: Column(
          children: [
            CircularProgressIndicator(
              color: _primaryColor,
              strokeWidth: 2,
            ),
            const SizedBox(height: 12),
            Text(
              'Loading payment history...',
              style: TextStyle(
                fontSize: 14,
                color: _textSecondary,
              ),
            ),
          ],
        ),
      );
    }

    if (_payments.isEmpty) {
      return Center(
        child: Column(
          children: [
            Icon(
              Icons.history_rounded,
              size: 48,
              color: _textSecondary,
            ),
            const SizedBox(height: 12),
            Text(
              'No payment history yet',
              style: TextStyle(
                fontSize: 16,
                color: _textSecondary,
              ),
            ),
          ],
        ),
      );
    }

    // Use the backend's calculated values
    List<Map<String, dynamic>> paymentsWithTotals = [];
    
    for (int i = 0; i < _payments.length; i++) {
      final payment = _payments[i];
      final newAddedRent = (payment['new_added_rent'] as num?)?.toDouble() ?? 0.0;
      final rent = (payment['rent'] as num?)?.toDouble() ?? 0.0;
      final receivedMoney = (payment['received_money'] as num?)?.toDouble() ?? 0.0;
      final dueRent = (payment['due_rent'] as num?)?.toDouble() ?? 0.0;
      final newAddedElectricityBill = (payment['new_added_electricity_bill'] as num?)?.toDouble();
      final paidElectricityBill = (payment['paid_electricity_bill'] as num?)?.toDouble();
      final dueElectricityBill = (payment['due_electricity_bill'] as num?)?.toDouble() ?? 0.0;
      final electricityBill = (payment['electricity_bill'] as num?)?.toDouble() ?? 0.0;
      
      paymentsWithTotals.add({
        ...payment,
        'new_added_rent': newAddedRent,
        'rent': rent,
        'due_rent': dueRent,
        'new_added_electricity_bill': newAddedElectricityBill,
        'paid_electricity_bill': paidElectricityBill,
        'due_electricity_bill': dueElectricityBill,
        'electricity_bill': electricityBill,
      });
    }

    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: DataTable(
        columnSpacing: 20,
        headingTextStyle: TextStyle(
          fontSize: 14,
          fontWeight: FontWeight.w700,
          color: _textPrimary,
        ),
        dataTextStyle: TextStyle(
          fontSize: 12,
          color: _textPrimary,
        ),
        columns: [
          DataColumn(label: Text(AppLocalizations.of(context).date)),
          DataColumn(label: Text(AppLocalizations.of(context).newAddedRent)),
          DataColumn(label: Text(AppLocalizations.of(context).rent)),
          DataColumn(label: Text(AppLocalizations.of(context).receivedMoney)),
          DataColumn(label: Text(AppLocalizations.of(context).dueRent)),
          DataColumn(label: Text(AppLocalizations.of(context).newAddedElectricityBill)),
          DataColumn(label: Text(AppLocalizations.of(context).electricityBill)),
          DataColumn(label: Text(AppLocalizations.of(context).paidElectricity)),
          DataColumn(label: Text(AppLocalizations.of(context).dueElectricityBill)),
        ],
        rows: paymentsWithTotals.map((payment) {
          final createdAt = payment['created_at'] as String? ?? '';
          final newAddedRent = payment['new_added_rent'] as double? ?? 0.0;
          final rent = (payment['rent'] as num?)?.toDouble() ?? 0.0;
          final receivedMoney = (payment['received_money'] as num?)?.toDouble() ?? 0.0;
          final dueRent = payment['due_rent'] as double? ?? 0.0;
          final fullPayment = payment['full_payment'] as bool? ?? false;
          final newAddedElectricityBill = payment['new_added_electricity_bill'] as double?;
          final electricityBill = payment['electricity_bill'] as double? ?? 0.0;
          final paidElectricityBill = payment['paid_electricity_bill'] as double?;
          final dueElectricityBill = payment['due_electricity_bill'] as double? ?? 0.0;

          return DataRow(
            cells: [
              DataCell(
                Text(
                  _formatDate(createdAt),
                  style: TextStyle(
                    fontSize: 11,
                    color: _textSecondary,
                  ),
                ),
              ),
              DataCell(
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: newAddedRent > 0 ? Colors.red.withOpacity(0.1) : Colors.green.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Text(
                    '${newAddedRent > 0 ? '+' : ''}${newAddedRent.toStringAsFixed(0)} ${AppLocalizations.of(context).tk}',
                    style: TextStyle(
                      fontSize: 11,
                      fontWeight: FontWeight.w600,
                      color: newAddedRent > 0 ? Colors.red.shade600 : Colors.green.shade600,
                    ),
                  ),
                ),
              ),
              DataCell(
                Text(
                  '${rent.toStringAsFixed(0)} ${AppLocalizations.of(context).tk}',
                  style: TextStyle(
                    fontSize: 11,
                    fontWeight: FontWeight.w600,
                    color: _textPrimary,
                  ),
                ),
              ),
              DataCell(
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: receivedMoney > 0 ? Colors.green.withOpacity(0.1) : Colors.grey.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Text(
                    '${receivedMoney.toStringAsFixed(0)} ${AppLocalizations.of(context).tk}',
                    style: TextStyle(
                      fontSize: 11,
                      fontWeight: FontWeight.w600,
                      color: receivedMoney > 0 ? Colors.green.shade600 : _textSecondary,
                    ),
                  ),
                ),
              ),
              DataCell(
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: dueRent > 0 ? Colors.red.withOpacity(0.1) : Colors.green.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Text(
                    '${dueRent.toStringAsFixed(0)} ${AppLocalizations.of(context).tk}',
                    style: TextStyle(
                      fontSize: 11,
                      fontWeight: FontWeight.w700,
                      color: dueRent > 0 ? Colors.red.shade600 : Colors.green.shade600,
                    ),
                  ),
                ),
              ),
              DataCell(
                newAddedElectricityBill != null
                    ? Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                        decoration: BoxDecoration(
                          color: newAddedElectricityBill > 0 ? Colors.orange.withOpacity(0.1) : Colors.green.withOpacity(0.1),
                          borderRadius: BorderRadius.circular(4),
                        ),
                        child: Text(
                          '${newAddedElectricityBill > 0 ? '+' : ''}${newAddedElectricityBill.toStringAsFixed(0)} ${AppLocalizations.of(context).tk}',
                          style: TextStyle(
                            fontSize: 11,
                            fontWeight: FontWeight.w600,
                            color: newAddedElectricityBill > 0 ? Colors.orange.shade600 : Colors.green.shade600,
                          ),
                        ),
                      )
                    : Text(
                        '-',
                        style: TextStyle(
                          fontSize: 11,
                          color: _textSecondary,
                        ),
                      ),
              ),
              DataCell(
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: electricityBill > 0 ? Colors.orange.withOpacity(0.1) : Colors.green.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Text(
                    '${electricityBill.toStringAsFixed(0)} ${AppLocalizations.of(context).tk}',
                    style: TextStyle(
                      fontSize: 11,
                      fontWeight: FontWeight.w600,
                      color: electricityBill > 0 ? Colors.orange.shade600 : Colors.green.shade600,
                    ),
                  ),
                ),
              ),
              DataCell(
                paidElectricityBill != null
                    ? Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                        decoration: BoxDecoration(
                          color: paidElectricityBill > 0 ? Colors.green.withOpacity(0.1) : Colors.grey.withOpacity(0.1),
                          borderRadius: BorderRadius.circular(4),
                        ),
                        child: Text(
                          '${paidElectricityBill.toStringAsFixed(0)} ${AppLocalizations.of(context).tk}',
                          style: TextStyle(
                            fontSize: 11,
                            fontWeight: FontWeight.w600,
                            color: paidElectricityBill > 0 ? Colors.green.shade600 : _textSecondary,
                          ),
                        ),
                      )
                    : Text(
                        '-',
                        style: TextStyle(
                          fontSize: 11,
                          color: _textSecondary,
                        ),
                      ),
              ),
              DataCell(
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: dueElectricityBill > 0 ? Colors.orange.withOpacity(0.1) : Colors.green.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Text(
                    '${dueElectricityBill.toStringAsFixed(0)} ${AppLocalizations.of(context).tk}',
                    style: TextStyle(
                      fontSize: 11,
                      fontWeight: FontWeight.w700,
                      color: dueElectricityBill > 0 ? Colors.orange.shade600 : Colors.green.shade600,
                    ),
                  ),
                ),
              ),
            ],
          );
        }).toList(),
      ),
    );
  }

  Widget _buildPaginationControls() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        // Previous button
        IconButton(
          onPressed: _hasPrevPage ? () => _loadPage(_currentPage - 1) : null,
          icon: Icon(
            Icons.chevron_left_rounded,
            color: _hasPrevPage ? _primaryColor : _textSecondary,
          ),
        ),
        
        // Page info
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
          decoration: BoxDecoration(
            color: _primaryColor.withOpacity(0.1),
            borderRadius: BorderRadius.circular(8),
          ),
          child: Text(
            'Page $_currentPage of $_totalPages',
            style: TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w600,
              color: _primaryColor,
            ),
          ),
        ),
        
        // Next button
        IconButton(
          onPressed: _hasNextPage ? () => _loadPage(_currentPage + 1) : null,
          icon: Icon(
            Icons.chevron_right_rounded,
            color: _hasNextPage ? _primaryColor : _textSecondary,
          ),
        ),
      ],
    );
  }

  Widget _buildPaymentDetail(String label, String value, Color color, IconData icon) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            Icon(icon, size: 16, color: color),
            const SizedBox(width: 4),
            Text(
              label,
              style: TextStyle(
                fontSize: 12,
                color: _textSecondary,
              ),
            ),
          ],
        ),
        const SizedBox(height: 4),
        Text(
          value,
          style: TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.w700,
            color: color,
          ),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: _backgroundColor,
      appBar: AppBar(
        title: Text(
          AppLocalizations.of(context).paymentDetails,
          style: TextStyle(
            fontWeight: FontWeight.w700,
            fontSize: 20,
            color: Colors.white,
          ),
        ),
        centerTitle: true,
        elevation: 0,
        backgroundColor: _primaryColor,
        systemOverlayStyle: SystemUiOverlayStyle.light,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back_rounded, color: Colors.white),
          onPressed: () => Navigator.pop(context),
        ),
      ),
      body: _isLoading
          ? const Center(
              child: CircularProgressIndicator(),
            )
          : _error != null
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Icon(
                        Icons.error_outline_rounded,
                        size: 64,
                        color: Colors.red.shade400,
                      ),
                      const SizedBox(height: 16),
                      Text(
                        'Error loading payment details',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.w600,
                          color: _textPrimary,
                        ),
                      ),
                      const SizedBox(height: 8),
                      Text(
                        _error!,
                        style: TextStyle(
                          color: _textSecondary,
                        ),
                        textAlign: TextAlign.center,
                      ),
                      const SizedBox(height: 16),
                      ElevatedButton(
                        onPressed: _loadPaymentDetails,
                        child: const Text('Try Again'),
                      ),
                    ],
                  ),
                )
              : SingleChildScrollView(
                  padding: const EdgeInsets.all(20),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [


                      // Payment Details Card
                      Container(
                        width: double.infinity,
                        padding: const EdgeInsets.all(20),
                        decoration: BoxDecoration(
                          color: _cardColor,
                          borderRadius: BorderRadius.circular(16),
                          boxShadow: [
                            BoxShadow(
                              color: Colors.black.withOpacity(0.05),
                              blurRadius: 10,
                              spreadRadius: 0,
                              offset: const Offset(0, 4),
                            ),
                          ],
                        ),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              AppLocalizations.of(context).paymentDetails,
                              style: TextStyle(
                                fontSize: 20,
                                fontWeight: FontWeight.w700,
                                color: _textPrimary,
                              ),
                            ),
                            const SizedBox(height: 20),
                            
                            // Rent (with adjust button for managers)
                            Row(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Expanded(
                                  child: Row(
                                    children: [
                                      Container(
                                        padding: const EdgeInsets.all(8),
                                        decoration: BoxDecoration(
                                          color: _primaryColor.withOpacity(0.1),
                                          borderRadius: BorderRadius.circular(8),
                                        ),
                                        child: Icon(
                                          Icons.attach_money_rounded,
                                          color: _primaryColor,
                                          size: 20,
                                        ),
                                      ),
                                      const SizedBox(width: 12),
                                      Expanded(
                                        child: Column(
                                          crossAxisAlignment: CrossAxisAlignment.start,
                                          children: [
                                            Text(
                                              AppLocalizations.of(context).rent,
                                              style: TextStyle(
                                                fontSize: 14,
                                                color: _textSecondary,
                                              ),
                                              overflow: TextOverflow.ellipsis,
                                            ),
                                            Text(
                                              '$_dueRent ${AppLocalizations.of(context).tk}',
                                              style: TextStyle(
                                                fontSize: 18,
                                                fontWeight: FontWeight.w700,
                                                color: _primaryColor,
                                              ),
                                              overflow: TextOverflow.ellipsis,
                                            ),
                                          ],
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                                if (_isManager)
                                  IconButton(
                                    onPressed: _showCreatePaymentDialog,
                                    icon: Icon(
                                      Icons.edit_rounded,
                                      color: _primaryColor,
                                      size: 28,
                                    ),
                                    tooltip: 'Adjust Due Rent',
                                  ),
                              ],
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(height: 20),

                      // Payment History Section
                      Container(
                        width: double.infinity,
                        padding: const EdgeInsets.all(20),
                        decoration: BoxDecoration(
                          color: _cardColor,
                          borderRadius: BorderRadius.circular(16),
                          boxShadow: [
                            BoxShadow(
                              color: Colors.black.withOpacity(0.05),
                              blurRadius: 10,
                              spreadRadius: 0,
                              offset: const Offset(0, 4),
                            ),
                          ],
                        ),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Container(
                                  padding: const EdgeInsets.all(8),
                                  decoration: BoxDecoration(
                                    color: _primaryColor.withOpacity(0.1),
                                    borderRadius: BorderRadius.circular(8),
                                  ),
                                  child: Icon(
                                    Icons.history_rounded,
                                    color: _primaryColor,
                                    size: 20,
                                  ),
                                ),
                                const SizedBox(width: 12),
                                Text(
                                  AppLocalizations.of(context).paymentHistory,
                                  style: TextStyle(
                                    fontSize: 20,
                                    fontWeight: FontWeight.w700,
                                    color: _textPrimary,
                                  ),
                                ),
                              ],
                            ),
                            const SizedBox(height: 16),
                            _buildPaymentHistoryTable(),
                            if (_totalPages > 1) ...[
                              const SizedBox(height: 16),
                              _buildPaginationControls(),
                            ],
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
    );
  }
} 