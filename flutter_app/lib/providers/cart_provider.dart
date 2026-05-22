import 'package:flutter/material.dart';
import '../models/cart_item.dart';
import '../services/api_service.dart';

class CartProvider extends ChangeNotifier {
  List<CartItem> _items = [];
  List<CartItem> get items => _items;

  int get itemCount => _items.length;

  double get total =>
      _items.fold(0, (sum, item) => sum + (item.price * item.quantity));

  Future<void> loadCart() async {
    final data = await ApiService.getCart();
    _items = data.map((e) => CartItem.fromJson(e)).toList();
    notifyListeners();
  }

  Future<bool> addToCart(String productId, int quantity) async {
    final success = await ApiService.addToCart(productId, quantity);
    if (success) await loadCart();
    return success;
  }

  Future<bool> checkout() async {
    final success = await ApiService.createOrder();
    if (success) {
      _items = [];
      notifyListeners();
    }
    return success;
  }
}
