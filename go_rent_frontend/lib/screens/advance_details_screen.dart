import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:intl/intl.dart';
import '../services/api_service.dart';
import '../utils/app_localizations.dart';

class AdvanceDetailsScreen extends StatefulWidget {
  final String floorName;
  final int floorId;
  final Color themeColor;

  const AdvanceDetailsScreen({
    Key? key,
    required this.floorName,
    required this.floorId,
    required this.themeColor,
  }) : super(key: key);

  @override
  _AdvanceDetailsScreenState createState() => _AdvanceDetailsScreenState();
}

class _AdvanceDetailsScreenState extends State<AdvanceDetailsScreen> {
  final ApiService _apiService = ApiService();
  List<Map<String, dynamic>> _advances = [];
  bool _isLoading = true;
  String? _error;
  bool _isManager = false;

  @override
  void initState() {
    super.initState();
    _loadAdvanceDetails();
    _checkManagerStatus();
  }

  Future<void> _loadAdvanceDetails() async {
    try {
      setState(() {
        _isLoading = true;
        _error = null;
      });

      final advances = await _apiService.getAdvanceDetails(widget.floorId);
      
      setState(() {
        _advances = advances;
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
    }
  }

  String _formatDate(String dateString) {
    try {
      final date = DateTime.parse(dateString);
      return DateFormat('MMM dd, yyyy').format(date);
    } catch (e) {
      return dateString;
    }
  }

  Future<void> _checkManagerStatus() async {
    try {
      // For now, we'll assume the user is a manager if they can access this screen
      // In a real implementation, you might want to pass the property ID from the parent screen
      setState(() {
        _isManager = true; // TODO: Implement proper manager check with property ID
      });
    } catch (e) {
      print('Error checking manager status: $e');
    }
  }

  void _showDeductAdvanceDialog(Map<String, dynamic> advance) {
    final TextEditingController amountController = TextEditingController();
    
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('Deduct Advance'),
          content: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Text(
                  'Current Advance: ${advance['money']} ${AppLocalizations.of(context).currency}',
                  style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
                ),
                const SizedBox(height: 16),
                TextField(
                  controller: amountController,
                  keyboardType: TextInputType.number,
                  decoration: InputDecoration(
                    labelText: 'Deduct Amount',
                    hintText: 'Enter amount to deduct',
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                ),
              ],
            ),
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () async {
                final amount = int.tryParse(amountController.text);
                
                if (amount == null || amount <= 0) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('Please enter a valid amount')),
                  );
                  return;
                }
                
                if (amount > advance['money']) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('Deduct amount cannot exceed advance amount')),
                  );
                  return;
                }
                
                Navigator.of(context).pop();
                await _deductAdvance(advance['id'], amount);
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: widget.themeColor,
                foregroundColor: Colors.white,
              ),
              child: const Text('Deduct'),
            ),
          ],
        );
      },
    );
  }

  Future<void> _deductAdvance(int advanceId, int amount) async {
    try {
      // Find the advance in the current list and update it
      final advanceIndex = _advances.indexWhere((advance) => advance['id'] == advanceId);
      if (advanceIndex != -1) {
        final currentAdvance = _advances[advanceIndex];
        final currentAmount = currentAdvance['money'] as int;
        final newAmount = currentAmount - amount;
        
        // Update the advance amount in the local list
        setState(() {
          _advances[advanceIndex] = {
            ...currentAdvance,
            'money': newAmount,
          };
        });
        
        // TODO: Implement the API call to update the advance in the database
        // await _apiService.updateAdvanceAmount(advanceId, newAmount);
        
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Successfully deducted $amount ${AppLocalizations.of(context).currency}'),
            backgroundColor: Colors.green,
          ),
        );
      } else {
        throw Exception('Advance not found');
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Error deducting advance: $e'),
          backgroundColor: Colors.red,
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.grey[50],
      appBar: AppBar(
        title: Text(
          AppLocalizations.of(context).advanceDetailsForFloor.replaceAll('{floorName}', widget.floorName),
          style: const TextStyle(
            fontWeight: FontWeight.w700,
            fontSize: 18,
            color: Colors.white,
          ),
        ),
        centerTitle: true,
        elevation: 0,
        backgroundColor: widget.themeColor,
        systemOverlayStyle: SystemUiOverlayStyle.light,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Colors.white),
          onPressed: () => Navigator.pop(context),
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh, color: Colors.white),
            onPressed: _loadAdvanceDetails,
          ),
        ],
      ),
      body: SafeArea(
        child: _isLoading
            ? const Center(child: CircularProgressIndicator())
            : _error != null
                ? Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(
                          Icons.error_outline,
                          size: 64,
                          color: Colors.grey[400],
                        ),
                        const SizedBox(height: 16),
                        Text(
                          _error!,
                          style: TextStyle(
                            fontSize: 16,
                            color: Colors.grey[600],
                          ),
                          textAlign: TextAlign.center,
                        ),
                        const SizedBox(height: 16),
                        ElevatedButton(
                          onPressed: _loadAdvanceDetails,
                          style: ElevatedButton.styleFrom(
                            backgroundColor: widget.themeColor,
                            foregroundColor: Colors.white,
                            padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(8),
                            ),
                          ),
                          child: Text(AppLocalizations.of(context).retry),
                        ),
                      ],
                    ),
                  )
                : _advances.isEmpty
                    ? Center(
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Icon(
                              Icons.account_balance_wallet_outlined,
                              size: 64,
                              color: Colors.grey[400],
                            ),
                            const SizedBox(height: 16),
                            Text(
                              AppLocalizations.of(context).noAdvanceDetails,
                              style: TextStyle(
                                fontSize: 16,
                                color: Colors.grey[600],
                              ),
                              textAlign: TextAlign.center,
                            ),
                          ],
                        ),
                      )
                    : RefreshIndicator(
                        onRefresh: _loadAdvanceDetails,
                        child: ListView.builder(
                          padding: const EdgeInsets.all(16),
                          itemCount: _advances.length,
                          itemBuilder: (context, index) {
                            final advance = _advances[index];
                            return Container(
                              margin: const EdgeInsets.only(bottom: 12),
                              decoration: BoxDecoration(
                                color: Colors.white,
                                borderRadius: BorderRadius.circular(12),
                                boxShadow: [
                                  BoxShadow(
                                    color: Colors.black.withOpacity(0.05),
                                    blurRadius: 8,
                                    offset: const Offset(0, 2),
                                  ),
                                ],
                              ),
                              child: Padding(
                                padding: const EdgeInsets.all(16),
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Row(
                                      children: [
                                        Container(
                                          padding: const EdgeInsets.all(8),
                                          decoration: BoxDecoration(
                                            color: widget.themeColor.withOpacity(0.1),
                                            borderRadius: BorderRadius.circular(8),
                                          ),
                                          child: Icon(
                                            Icons.account_balance_wallet_rounded,
                                            color: widget.themeColor,
                                            size: 20,
                                          ),
                                        ),
                                        const SizedBox(width: 12),
                                        Expanded(
                                          child: Text(
                                            advance['user_name'] ?? 'Unknown User',
                                            style: const TextStyle(
                                              fontSize: 16,
                                              fontWeight: FontWeight.w600,
                                              color: Colors.black87,
                                            ),
                                          ),
                                        ),
                                        Container(
                                          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                                          decoration: BoxDecoration(
                                            color: (advance['status'] == 'accepted' ? Colors.green : Colors.orange).withOpacity(0.1),
                                            borderRadius: BorderRadius.circular(6),
                                          ),
                                          child: Text(
                                            advance['status'] ?? 'Unknown',
                                            style: TextStyle(
                                              fontSize: 12,
                                              fontWeight: FontWeight.w500,
                                              color: advance['status'] == 'accepted' ? Colors.green : Colors.orange,
                                            ),
                                          ),
                                        ),
                                      ],
                                    ),
                                    const SizedBox(height: 16),
                                    Row(
                                      children: [
                                        Expanded(
                                          child: Column(
                                            crossAxisAlignment: CrossAxisAlignment.start,
                                            children: [
                                              Text(
                                                AppLocalizations.of(context).amount,
                                                style: TextStyle(
                                                  fontSize: 12,
                                                  color: Colors.grey[600],
                                                ),
                                              ),
                                              const SizedBox(height: 4),
                                              Text(
                                                '${advance['money']} ${AppLocalizations.of(context).currency}',
                                                style: const TextStyle(
                                                  fontSize: 18,
                                                  fontWeight: FontWeight.w700,
                                                  color: Colors.black87,
                                                ),
                                              ),
                                            ],
                                          ),
                                        ),
                                        Expanded(
                                          child: Column(
                                            crossAxisAlignment: CrossAxisAlignment.start,
                                            children: [
                                              Text(
                                                AppLocalizations.of(context).date,
                                                style: TextStyle(
                                                  fontSize: 12,
                                                  color: Colors.grey[600],
                                                ),
                                              ),
                                              const SizedBox(height: 4),
                                              Text(
                                                _formatDate(advance['created_at'] ?? ''),
                                                style: const TextStyle(
                                                  fontSize: 14,
                                                  fontWeight: FontWeight.w500,
                                                  color: Colors.black87,
                                                ),
                                              ),
                                            ],
                                          ),
                                        ),
                                      ],
                                    ),
                                    // Deduct Advance button for managers
                                    if (_isManager && advance['status'] == 'accepted') ...[
                                      const SizedBox(height: 16),
                                      SizedBox(
                                        width: double.infinity,
                                        child: ElevatedButton.icon(
                                          icon: const Icon(Icons.remove_circle_outline, size: 16),
                                          label: const Text('Deduct Advance'),
                                          style: ElevatedButton.styleFrom(
                                            backgroundColor: Colors.red,
                                            foregroundColor: Colors.white,
                                            padding: const EdgeInsets.symmetric(vertical: 12),
                                            shape: RoundedRectangleBorder(
                                              borderRadius: BorderRadius.circular(8),
                                            ),
                                          ),
                                          onPressed: () => _showDeductAdvanceDialog(advance),
                                        ),
                                      ),
                                    ],
                                  ],
                                ),
                              ),
                            );
                          },
                        ),
                      ),
      ),
    );
  }
} 