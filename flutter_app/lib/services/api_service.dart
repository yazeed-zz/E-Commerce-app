import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';

class ApiService {
  static const String baseUrl = 'http://192.168.0.199:8080/api';

  static Future<String?> getToken() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString('token');
  }

  static Future<Map<String, String>> getHeaders() async {
    final token = await getToken();
    return {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $token',
    };
  }

  // Products
  static Future<List<dynamic>> getProducts({
    String? search,
    String? category,
  }) async {
    final headers = await getHeaders();
    var url = '$baseUrl/products';
    if (search != null) url += '?search=$search';
    if (category != null) url += '?category=$category';

    final response = await http.get(Uri.parse(url), headers: headers);
    final data = jsonDecode(response.body);
    return data['products'] ?? [];
  }

  static Future<Map<String, dynamic>> getProduct(String id) async {
    final headers = await getHeaders();
    final response = await http.get(
      Uri.parse('$baseUrl/products/$id'),
      headers: headers,
    );
    return jsonDecode(response.body);
  }

  // Recommendations
  static Future<List<dynamic>> getRecommendations(
    String type, {
    String? id,
  }) async {
    final headers = await getHeaders();
    var url = '$baseUrl/recommendations/$type';
    if (id != null) url += '?id=$id';

    final response = await http.get(Uri.parse(url), headers: headers);
    final data = jsonDecode(response.body);
    return data['products'] ?? [];
  }

  // Cart
  static Future<bool> addToCart(String productId, int quantity) async {
    final headers = await getHeaders();
    final response = await http.post(
      Uri.parse('$baseUrl/cart'),
      headers: headers,
      body: jsonEncode({'product_id': productId, 'quantity': quantity}),
    );
    return response.statusCode == 200;
  }

  static Future<List<dynamic>> getCart() async {
    final headers = await getHeaders();
    final response = await http.get(
      Uri.parse('$baseUrl/cart'),
      headers: headers,
    );
    final data = jsonDecode(response.body);
    return data['cart'] ?? [];
  }

  // Orders
  static Future<bool> createOrder() async {
    final headers = await getHeaders();
    final response = await http.post(
      Uri.parse('$baseUrl/orders'),
      headers: headers,
    );
    return response.statusCode == 201;
  }

  static Future<List<dynamic>> getOrders() async {
    final headers = await getHeaders();
    final response = await http.get(
      Uri.parse('$baseUrl/orders'),
      headers: headers,
    );
    final data = jsonDecode(response.body);
    return data['orders'] ?? [];
  }

  // Admin
  static Future<Map<String, dynamic>> getStats() async {
    final headers = await getHeaders();
    final response = await http.get(
      Uri.parse('$baseUrl/admin/stats'),
      headers: headers,
    );
    return jsonDecode(response.body);
  }

  static Future<List<dynamic>> getUsers() async {
    final headers = await getHeaders();
    final response = await http.get(
      Uri.parse('$baseUrl/admin/users'),
      headers: headers,
    );
    final data = jsonDecode(response.body);
    return data['users'] ?? [];
  }

  static Future<bool> createProduct(Map<String, dynamic> product) async {
    final headers = await getHeaders();
    final response = await http.post(
      Uri.parse('$baseUrl/admin/products'),
      headers: headers,
      body: jsonEncode(product),
    );
    return response.statusCode == 201;
  }

  static Future<bool> deleteProduct(String id) async {
    final headers = await getHeaders();
    final response = await http.delete(
      Uri.parse('$baseUrl/admin/products/$id'),
      headers: headers,
    );
    return response.statusCode == 200;
  }

  static Future<bool> updateOrderStatus(String id, String status) async {
    final headers = await getHeaders();
    final response = await http.put(
      Uri.parse('$baseUrl/admin/orders/$id/status'),
      headers: headers,
      body: jsonEncode({'status': status}),
    );
    return response.statusCode == 200;
  }
}
