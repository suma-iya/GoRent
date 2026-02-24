import 'package:flutter/material.dart';
import 'package:flutter_markdown/flutter_markdown.dart';
import 'package:provider/provider.dart';
import 'package:url_launcher/url_launcher.dart';
import '../providers/chat_provider.dart';

class ChatScreen extends StatefulWidget {
  const ChatScreen({super.key});

  @override
  State<ChatScreen> createState() => _ChatScreenState();
}

class _ChatScreenState extends State<ChatScreen> {
  final TextEditingController _controller = TextEditingController();
  final ScrollController _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _scrollToBottom();
    });
  }

  @override
  void dispose() {
    _controller.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  void _scrollToBottom() {
    if (_scrollController.hasClients) {
      Future.delayed(const Duration(milliseconds: 100), () {
        if (_scrollController.hasClients) {
          _scrollController.animateTo(
            _scrollController.position.maxScrollExtent,
            duration: const Duration(milliseconds: 300),
            curve: Curves.easeOut,
          );
        }
      });
    }
  }

  void _sendMessage() {
    final text = _controller.text.trim();
    if (text.isEmpty) return;

    _controller.clear();
    FocusScope.of(context).unfocus();

    final chatProvider = Provider.of<ChatProvider>(context, listen: false);
    chatProvider.sendMessage(text);
    
    // Scroll to bottom after sending
    Future.delayed(const Duration(milliseconds: 100), () {
      _scrollToBottom();
    });
  }

  Widget _buildMessageText(ChatMessage message, BuildContext context) {
    // Clean the text to remove any problematic characters
    String cleanText = message.text
        .replaceAll('â¢', '-')  // Fix encoding issues
        .replaceAll('â€¢', '-')  // Fix other encoding variants
        .replaceAll('\u2022', '-')  // Unicode bullet to dash
        .replaceAll('\u00A0', ' ');  // Non-breaking space to regular space
    
    try {
      return MarkdownBody(
        data: cleanText,
        styleSheet: MarkdownStyleSheet(
          p: TextStyle(
            fontSize: 14,
            color: message.isUser 
                ? Theme.of(context).primaryColor
                : Colors.grey[900],
            height: 1.4,
          ),
          strong: TextStyle(
            fontWeight: FontWeight.bold,
            color: message.isUser 
                ? Theme.of(context).primaryColor
                : Colors.grey[900],
          ),
          listBullet: TextStyle(
            color: message.isUser 
                ? Theme.of(context).primaryColor
                : Colors.grey[900],
          ),
          code: TextStyle(
            backgroundColor: Colors.grey[200],
            fontFamily: 'monospace',
          ),
          h1: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: message.isUser 
                ? Theme.of(context).primaryColor
                : Colors.grey[900],
          ),
          h2: TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.bold,
            color: message.isUser 
                ? Theme.of(context).primaryColor
                : Colors.grey[900],
          ),
        ),
        onTapLink: (text, href, title) {
          if (href != null) {
            launchUrl(Uri.parse(href));
          }
        },
        shrinkWrap: true,
      );
    } catch (e) {
      print('Markdown rendering error: $e');
      // Fallback to plain text if markdown fails
      return SelectableText(
        cleanText,
        style: TextStyle(
          fontSize: 14,
          color: message.isUser 
              ? Theme.of(context).primaryColor
              : Colors.grey[900],
          height: 1.4,
        ),
      );
    }
  }

  Widget _buildMessage(ChatMessage message) {
    return Container(
      margin: const EdgeInsets.symmetric(vertical: 4, horizontal: 8),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisAlignment: message.isUser ? MainAxisAlignment.end : MainAxisAlignment.start,
        children: [
          if (!message.isUser)
            const CircleAvatar(
              backgroundColor: Colors.blue,
              radius: 16,
              child: Icon(Icons.smart_toy, color: Colors.white, size: 18),
            ),
          Expanded(
            child: Column(
              crossAxisAlignment: message.isUser ? CrossAxisAlignment.end : CrossAxisAlignment.start,
              children: [
                Container(
                  constraints: BoxConstraints(
                    maxWidth: MediaQuery.of(context).size.width * 0.75,
                  ),
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: message.isUser 
                        ? Theme.of(context).primaryColor.withOpacity(0.1)
                        : Colors.grey[100],
                    borderRadius: BorderRadius.circular(16),
                    border: Border.all(
                      color: message.isUser 
                          ? Theme.of(context).primaryColor.withOpacity(0.3)
                          : Colors.grey[300]!,
                    ),
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      // Use markdown if available, otherwise plain text
                      message.text.isNotEmpty
                          ? _buildMessageText(message, context)
                          : Text(
                              'Empty message',
                              style: TextStyle(
                                fontSize: 14,
                                color: Colors.grey[500],
                                fontStyle: FontStyle.italic,
                              ),
                            ),
                      if (!message.isUser && message.intent != null)
                        Padding(
                          padding: const EdgeInsets.only(top: 8),
                          child: Row(
                            children: [
                              Container(
                                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                                decoration: BoxDecoration(
                                  color: Colors.blue[100],
                                  borderRadius: BorderRadius.circular(12),
                                ),
                                child: Text(
                                  message.intent!.replaceAll('_', ' '),
                                  style: const TextStyle(
                                    fontSize: 10,
                                    color: Colors.blue,
                                    fontWeight: FontWeight.w500,
                                  ),
                                ),
                              ),
                              if (message.confidence != null)
                                Padding(
                                  padding: const EdgeInsets.only(left: 8),
                                  child: Text(
                                    '${(message.confidence! * 100).toStringAsFixed(0)}%',
                                    style: const TextStyle(
                                      fontSize: 10,
                                      color: Colors.grey,
                                    ),
                                  ),
                                ),
                            ],
                          ),
                        ),
                    ],
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.only(top: 4, left: 8, right: 8),
                  child: Text(
                    message.timestamp,
                    style: const TextStyle(
                      fontSize: 10,
                      color: Colors.grey,
                    ),
                  ),
                ),
                if (!message.isUser && message.suggestedFollowups != null && message.suggestedFollowups!.isNotEmpty)
                  Padding(
                    padding: const EdgeInsets.only(top: 8),
                    child: Wrap(
                      spacing: 8,
                      runSpacing: 4,
                      children: message.suggestedFollowups!.take(3).map((followup) {
                        return ActionChip(
                          label: Text(followup),
                          onPressed: () {
                            final chatProvider = Provider.of<ChatProvider>(context, listen: false);
                            _controller.text = followup;
                            _sendMessage();
                          },
                          backgroundColor: Colors.blue[50],
                          labelStyle: const TextStyle(fontSize: 12),
                        );
                      }).toList(),
                    ),
                  ),
              ],
            ),
          ),
          if (message.isUser)
            const CircleAvatar(
              backgroundColor: Colors.green,
              radius: 16,
              child: Icon(Icons.person, color: Colors.white, size: 18),
            ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final chatProvider = Provider.of<ChatProvider>(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Rent Decision Chatbot'),
        actions: [
          if (chatProvider.selectedTenantId != null)
            Container(
              margin: const EdgeInsets.only(right: 8),
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
              decoration: BoxDecoration(
                color: Colors.blue[50],
                borderRadius: BorderRadius.circular(20),
                border: Border.all(color: Colors.blue[100]!),
              ),
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  const Icon(Icons.person, size: 16, color: Colors.blue),
                  const SizedBox(width: 4),
                  Text(
                    chatProvider.selectedTenantId!,
                    style: const TextStyle(
                      color: Colors.blue,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                ],
              ),
            ),
          IconButton(
            icon: const Icon(Icons.delete_outline),
            onPressed: () {
              showDialog(
                context: context,
                builder: (context) => AlertDialog(
                  title: const Text('Clear Chat'),
                  content: const Text('Are you sure you want to clear all chat messages?'),
                  actions: [
                    TextButton(
                      onPressed: () => Navigator.pop(context),
                      child: const Text('Cancel'),
                    ),
                    TextButton(
                      onPressed: () {
                        Navigator.pop(context);
                        chatProvider.clearHistory();
                      },
                      child: const Text('Clear'),
                    ),
                  ],
                ),
              );
            },
            tooltip: 'Clear chat',
          ),
          IconButton(
            icon: const Icon(Icons.help_outline),
            onPressed: () {
              showDialog(
                context: context,
                builder: (context) => AlertDialog(
                  title: const Text('How to use'),
                  content: const SingleChildScrollView(
                    child: Text(
                      '• Ask about specific tenants (e.g., "Why is tenant with phone 01712345678 high risk?")\n\n'
                      '• Get recommendations (e.g., "What should I do for tenant 01987654321?")\n\n'
                      '• View summaries (e.g., "Show monthly summary")\n\n'
                      '• Compare tenants (e.g., "Compare tenants 01712345678 and 01987654321")\n\n'
                      '• List high-risk tenants (e.g., "List high risk tenants")',
                    ),
                  ),
                  actions: [
                    TextButton(
                      onPressed: () => Navigator.pop(context),
                      child: const Text('OK'),
                    ),
                  ],
                ),
              );
            },
            tooltip: 'Help',
          ),
        ],
      ),
      body: Column(
        children: [
          // Quick actions bar
          Container(
            padding: const EdgeInsets.symmetric(vertical: 8, horizontal: 12),
            color: Colors.grey[50],
            child: SingleChildScrollView(
              scrollDirection: Axis.horizontal,
              child: Row(
                children: [
                  _buildQuickAction(
                    context,
                    Icons.warning,
                    'High Risk',
                    'List high risk tenants',
                  ),
                  _buildQuickAction(
                    context,
                    Icons.summarize,
                    'Summary',
                    'Show monthly summary',
                  ),
                  _buildQuickAction(
                    context,
                    Icons.compare,
                    'Compare',
                    'Compare tenants 01712345678 and 01987654321',
                  ),
                  _buildQuickAction(
                    context,
                    Icons.assignment,
                    'Actions',
                    chatProvider.selectedTenantId != null
                        ? 'What should I do for tenant ${chatProvider.selectedTenantId}?'
                        : 'What should I do for tenant 01712345678?',
                  ),
                ],
              ),
            ),
          ),

          // Messages list
          Expanded(
            child: chatProvider.messages.isEmpty
                ? Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(
                          Icons.chat_bubble_outline,
                          size: 64,
                          color: Colors.grey[400],
                        ),
                        const SizedBox(height: 16),
                        Text(
                          'Start a conversation',
                          style: TextStyle(
                            fontSize: 18,
                            color: Colors.grey[600],
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          'Ask about tenant risks, recommendations, or summaries',
                          style: TextStyle(
                            fontSize: 14,
                            color: Colors.grey[500],
                          ),
                          textAlign: TextAlign.center,
                        ),
                      ],
                    ),
                  )
                : ListView.builder(
                    controller: _scrollController,
                    itemCount: chatProvider.messages.length,
                    padding: const EdgeInsets.only(bottom: 8),
                    itemBuilder: (context, index) {
                      if (index >= chatProvider.messages.length) {
                        return const SizedBox.shrink();
                      }
                      try {
                        return _buildMessage(chatProvider.messages[index]);
                      } catch (e) {
                        print('Error building message at index $index: $e');
                        return Container(
                          padding: const EdgeInsets.all(8),
                          margin: const EdgeInsets.symmetric(vertical: 4, horizontal: 8),
                          child: Text(
                            'Error displaying message: ${chatProvider.messages[index].text.substring(0, chatProvider.messages[index].text.length > 50 ? 50 : chatProvider.messages[index].text.length)}',
                            style: const TextStyle(color: Colors.red),
                          ),
                        );
                      }
                    },
                  ),
          ),

          // Loading indicator
          if (chatProvider.isLoading)
            const LinearProgressIndicator(
              minHeight: 2,
              backgroundColor: Colors.transparent,
            ),

          // Input area
          Container(
            padding: const EdgeInsets.all(8),
            decoration: BoxDecoration(
              color: Theme.of(context).scaffoldBackgroundColor,
              border: Border(top: BorderSide(color: Colors.grey[300]!)),
              boxShadow: [
                BoxShadow(
                  color: Colors.grey.withOpacity(0.1),
                  blurRadius: 10,
                  offset: const Offset(0, -5),
                ),
              ],
            ),
            child: Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _controller,
                    decoration: InputDecoration(
                      hintText: chatProvider.selectedTenantId != null
                          ? 'Ask about ${chatProvider.selectedTenantId}...'
                          : 'Ask about tenant risk...',
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(24),
                      ),
                      contentPadding: const EdgeInsets.symmetric(
                        vertical: 12,
                        horizontal: 16,
                      ),
                    ),
                    onSubmitted: (_) => _sendMessage(),
                    maxLines: 3,
                    minLines: 1,
                  ),
                ),
                const SizedBox(width: 8),
                IconButton(
                  icon: const Icon(Icons.send),
                  onPressed: _sendMessage,
                  color: Theme.of(context).primaryColor,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildQuickAction(BuildContext context, IconData icon, String label, String query) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 4),
      child: FilterChip(
        label: Text(label),
        avatar: Icon(icon, size: 18),
        onSelected: (_) {
          final chatProvider = Provider.of<ChatProvider>(context, listen: false);
          _controller.text = query;
          _sendMessage();
        },
      ),
    );
  }
}

