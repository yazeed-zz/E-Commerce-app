import 'package:flutter/material.dart';
import '../models/product.dart';

class ProductCard extends StatelessWidget {
  final Product product;
  final VoidCallback? onTap;

  const ProductCard({super.key, required this.product, this.onTap});

  static Color _getProductColor(String name) {
    if (name.toLowerCase().contains('nike') ||
        name.toLowerCase().contains('shoe')) {
      return Colors.deepOrange;
    } else if (name.toLowerCase().contains('apple') ||
        name.toLowerCase().contains('watch')) {
      return Colors.blueGrey;
    } else if (name.toLowerCase().contains('sony') ||
        name.toLowerCase().contains('headphone')) {
      return Colors.indigo;
    }
    return Colors.deepPurple;
  }

  static IconData _getProductIcon(String name) {
    if (name.toLowerCase().contains('nike') ||
        name.toLowerCase().contains('shoe')) {
      return Icons.directions_run;
    } else if (name.toLowerCase().contains('apple') ||
        name.toLowerCase().contains('watch')) {
      return Icons.watch;
    } else if (name.toLowerCase().contains('sony') ||
        name.toLowerCase().contains('headphone')) {
      return Icons.headphones;
    }
    return Icons.shopping_bag;
  }

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Card(
        elevation: 2,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Expanded(
              child: ClipRRect(
                borderRadius: const BorderRadius.vertical(
                  top: Radius.circular(12),
                ),
                child: product.images.isNotEmpty
                    ? Image.network(
                        product.images[0],
                        fit: BoxFit.cover,
                        width: double.infinity,
                        errorBuilder: (_, __, ___) => Container(
                          color: _getProductColor(product.name),
                          child: Center(
                            child: Icon(
                              _getProductIcon(product.name),
                              size: 60,
                              color: Colors.white,
                            ),
                          ),
                        ),
                        loadingBuilder: (context, child, loadingProgress) {
                          if (loadingProgress == null) return child;
                          return Container(
                            color: _getProductColor(product.name),
                            child: const Center(
                              child: CircularProgressIndicator(
                                color: Colors.white,
                              ),
                            ),
                          );
                        },
                      )
                    : Container(
                        color: _getProductColor(product.name),
                        child: Center(
                          child: Icon(
                            _getProductIcon(product.name),
                            size: 60,
                            color: Colors.white,
                          ),
                        ),
                      ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(8),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    product.name,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                    style: const TextStyle(fontWeight: FontWeight.bold),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    '\$${product.price.toStringAsFixed(2)}',
                    style: const TextStyle(
                      color: Colors.deepPurple,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  Row(
                    children: [
                      const Icon(Icons.star, size: 14, color: Colors.amber),
                      Text(
                        product.ratingAvg.toStringAsFixed(1),
                        style: const TextStyle(fontSize: 12),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
