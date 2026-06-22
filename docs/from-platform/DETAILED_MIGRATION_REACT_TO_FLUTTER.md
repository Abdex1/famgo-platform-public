# 📖 DETAILED MIGRATION: React Components to Flutter Widgets

**Migration Guide**: Component-by-component conversion from React to Flutter  
**Focus**: Driver and User modules  
**Status**: Ready for implementation  

---

## 🎯 COMPONENT 1: RideBooking (React → Flutter)

### REACT VERSION (TypeScript)
```typescript
// C:\dev\FamGo\src\components\user\RideBooking\RideBooking.tsx
// Key Features:
// - Leaflet map for location selection
// - Real-time driver location display
// - Fare estimation
// - Ride pooling toggle
// - Join existing pool
```

### FLUTTER EQUIVALENT

**File**: `lib/features/user/presentation/screens/ride_booking_screen.dart`

```dart
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:geolocator/geolocator.dart';

// Import models and controllers
import '../../domain/models/ride.dart';
import '../controllers/ride_booking_controller.dart';
import '../widgets/location_search_widget.dart';
import '../widgets/fare_estimate_widget.dart';

class RideBookingScreen extends StatefulWidget {
  @override
  State<RideBookingScreen> createState() => _RideBookingScreenState();
}

class _RideBookingScreenState extends State<RideBookingScreen> {
  final controller = Get.put(RideBookingController());
  GoogleMapController? _mapController;
  final Set<Marker> _markers = {};
  final Set<Polyline> _polylines = {};

  @override
  void initState() {
    super.initState();
    _checkPoolJoinParam();
    controller.findNearbyDrivers();
  }

  void _checkPoolJoinParam() {
    final poolId = Get.parameters['poolId'];
    if (poolId != null) {
      controller.joinPool(poolId);
    }
  }

  void _updateMap() {
    _markers.clear();
    _polylines.clear();

    // Add pickup marker
    if (controller.pickupLocation.value != null) {
      _markers.add(
        Marker(
          markerId: MarkerId('pickup'),
          position: LatLng(
            controller.pickupLocation.value!.lat,
            controller.pickupLocation.value!.lng,
          ),
          infoWindow: InfoWindow(title: 'Pickup'),
          icon: BitmapDescriptor.defaultMarkerWithHue(
            BitmapDescriptor.hueGreen,
          ),
        ),
      );
    }

    // Add dropoff marker
    if (controller.dropoffLocation.value != null) {
      _markers.add(
        Marker(
          markerId: MarkerId('dropoff'),
          position: LatLng(
            controller.dropoffLocation.value!.lat,
            controller.dropoffLocation.value!.lng,
          ),
          infoWindow: InfoWindow(title: 'Dropoff'),
          icon: BitmapDescriptor.defaultMarkerWithHue(
            BitmapDescriptor.hueRed,
          ),
        ),
      );
    }

    // Add nearby driver markers
    int index = 0;
    for (var driver in controller.nearbyDrivers) {
      _markers.add(
        Marker(
          markerId: MarkerId('driver_$index'),
          position: LatLng(driver.location.lat, driver.location.lng),
          infoWindow: InfoWindow(
            title: driver.name,
            snippet: driver.vehicleType,
          ),
          icon: BitmapDescriptor.defaultMarkerWithHue(
            BitmapDescriptor.hueCyan,
          ),
        ),
      );
      index++;
    }

    // Draw polyline between pickup and dropoff
    if (controller.pickupLocation.value != null &&
        controller.dropoffLocation.value != null) {
      _polylines.add(
        Polyline(
          polylineId: PolylineId('route'),
          points: [
            LatLng(
              controller.pickupLocation.value!.lat,
              controller.pickupLocation.value!.lng,
            ),
            LatLng(
              controller.dropoffLocation.value!.lat,
              controller.dropoffLocation.value!.lng,
            ),
          ],
          color: Colors.blue,
          width: 4,
          patterns: [PatternItem.dash(20)],
        ),
      );
    }

    setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(
          controller.isJoiningPool.value ? 'Join a Pool' : 'Book a Ride',
        ),
      ),
      body: GetBuilder<RideBookingController>(
        builder: (_) {
          _updateMap();

          return SingleChildScrollView(
            child: Column(
              children: [
                // Pool info banner
                if (controller.isJoiningPool.value && controller.poolInfo.value != null)
                  _buildPoolInfoBanner(),

                // Map
                _buildMapWidget(),

                // Booking form
                Padding(
                  padding: EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Trip Details',
                        style: Theme.of(context).textTheme.titleLarge,
                      ),
                      SizedBox(height: 16),
                      
                      // Location inputs
                      if (controller.isJoiningPool.value)
                        _buildReadOnlyLocations()
                      else
                        _buildEditableLocations(),

                      SizedBox(height: 16),

                      // Pooling toggle
                      if (!controller.isJoiningPool.value)
                        _buildPoolingToggle(),

                      SizedBox(height: 16),

                      // Fare estimate
                      if (controller.fareInfo.value != null)
                        FareEstimateWidget(
                          fare: controller.fareInfo.value!,
                        ),

                      SizedBox(height: 16),

                      // Book/Join button
                      SizedBox(
                        width: double.infinity,
                        child: ElevatedButton(
                          onPressed: controller.pickupLocation.value != null &&
                                  controller.dropoffLocation.value != null
                              ? () => _handleBookRide()
                              : null,
                          child: Obx(
                            () => controller.isLoading.value
                                ? CircularProgressIndicator()
                                : Text(
                                    controller.isJoiningPool.value
                                        ? 'Join Pool'
                                        : 'Request Ride',
                                  ),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),

                // Nearby drivers list
                if (controller.nearbyDrivers.isNotEmpty)
                  _buildNearbyDriversList(),
              ],
            ),
          );
        },
      ),
    );
  }

  Widget _buildMapWidget() {
    return Container(
      height: 400,
      child: GoogleMap(
        onMapCreated: (controller) => _mapController = controller,
        initialCameraPosition: CameraPosition(
          target: LatLng(9.0320, 38.7469), // Addis Ababa
          zoom: 12,
        ),
        markers: _markers,
        polylines: _polylines,
        onTap: (latLng) {
          if (controller.isJoiningPool.value) return; // Can't change locations when joining

          if (controller.pickupLocation.value == null) {
            controller.setPickupLocation(latLng);
          } else if (controller.dropoffLocation.value == null) {
            controller.setDropoffLocation(latLng);
          }
        },
      ),
    );
  }

  Widget _buildReadOnlyLocations() {
    return Column(
      children: [
        _buildLocationDisplay(
          'Pickup Location',
          controller.pickupLocation.value?.address ?? 'Loading...',
          Colors.green,
        ),
        SizedBox(height: 12),
        _buildLocationDisplay(
          'Dropoff Location',
          controller.dropoffLocation.value?.address ?? 'Loading...',
          Colors.red,
        ),
        SizedBox(height: 12),
        Container(
          padding: EdgeInsets.all(12),
          decoration: BoxDecoration(
            color: Colors.orange.withOpacity(0.1),
            borderRadius: BorderRadius.circular(8),
          ),
          child: Row(
            children: [
              Icon(Icons.lock, color: Colors.orange, size: 20),
              SizedBox(width: 8),
              Text(
                'Locations locked for pool',
                style: TextStyle(color: Colors.orange),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildLocationDisplay(
    String label,
    String address,
    Color color,
  ) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: TextStyle(fontWeight: FontWeight.bold)),
        SizedBox(height: 8),
        Container(
          padding: EdgeInsets.all(12),
          decoration: BoxDecoration(
            border: Border.all(color: Colors.grey),
            borderRadius: BorderRadius.circular(8),
          ),
          child: Row(
            children: [
              Container(
                width: 12,
                height: 12,
                decoration: BoxDecoration(
                  color: color,
                  shape: BoxShape.circle,
                ),
              ),
              SizedBox(width: 12),
              Expanded(child: Text(address)),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildEditableLocations() {
    return Column(
      children: [
        LocationSearchWidget(
          label: 'Pickup Location',
          initialValue: controller.pickupLocation.value?.address,
          onChanged: (location) {
            controller.setPickupLocation(
              LatLng(location.lat, location.lng),
            );
            controller.findNearbyDrivers();
          },
        ),
        SizedBox(height: 12),
        LocationSearchWidget(
          label: 'Dropoff Location',
          initialValue: controller.dropoffLocation.value?.address,
          onChanged: (location) {
            controller.setDropoffLocation(
              LatLng(location.lat, location.lng),
            );
          },
        ),
      ],
    );
  }

  Widget _buildPoolingToggle() {
    return Container(
      padding: EdgeInsets.all(12),
      decoration: BoxDecoration(
        border: Border.all(color: Colors.grey.withOpacity(0.3)),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'Enable Pooling',
                style: TextStyle(fontWeight: FontWeight.bold),
              ),
              Text('Save up to 25% by sharing'),
            ],
          ),
          Obx(
            () => Switch(
              value: controller.wantPooling.value,
              onChanged: (value) => controller.setPooling(value),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildPoolInfoBanner() {
    final pool = controller.poolInfo.value!;
    return Container(
      margin: EdgeInsets.all(16),
      padding: EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.blue.withOpacity(0.1),
        border: Border.all(color: Colors.blue),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '🚙 Joining Pool',
            style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
          ),
          SizedBox(height: 8),
          Text('From: ${pool.pickupLocation.address}'),
          Text('To: ${pool.dropoffLocation.address}'),
          Text(
            'Passengers: ${pool.currentPassengers}/${pool.maxPassengers}',
          ),
        ],
      ),
    );
  }

  Widget _buildNearbyDriversList() {
    return Padding(
      padding: EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Nearby Drivers (${controller.nearbyDrivers.length})',
            style: Theme.of(context).textTheme.titleLarge,
          ),
          SizedBox(height: 12),
          ...controller.nearbyDrivers.map((driver) {
            return Card(
              child: ListTile(
                leading: CircleAvatar(
                  child: Text(driver.name[0]),
                ),
                title: Text(driver.name),
                subtitle: Text(
                  '${driver.vehicleType} • ${driver.distance.toStringAsFixed(1)} km',
                ),
                trailing: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(Icons.star, color: Colors.amber, size: 16),
                    SizedBox(width: 4),
                    Text(driver.rating.toStringAsFixed(1)),
                  ],
                ),
              ),
            );
          }).toList(),
        ],
      ),
    );
  }

  Future<void> _handleBookRide() async {
    try {
      if (controller.isJoiningPool.value) {
        await controller.joinExistingPool();
        Get.snackbar(
          'Success',
          'Successfully joined the pool!',
          backgroundColor: Colors.green,
        );
      } else {
        await controller.requestRide();
        Get.snackbar(
          'Success',
          'Ride requested successfully!',
          backgroundColor: Colors.green,
        );
      }
      Future.delayed(Duration(seconds: 2), () {
        Get.offNamed('/user/dashboard');
      });
    } catch (e) {
      Get.snackbar(
        'Error',
        'Failed to request ride: $e',
        backgroundColor: Colors.red,
      );
    }
  }

  @override
  void dispose() {
    _mapController?.dispose();
    super.dispose();
  }
}
```

---

## 🎯 COMPONENT 2: ActiveRide (React → Flutter)

**File**: `lib/features/driver/presentation/screens/active_ride_screen.dart`

```dart
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import '../../domain/models/ride.dart';
import '../controllers/active_ride_controller.dart';

class ActiveRideScreen extends StatefulWidget {
  @override
  State<ActiveRideScreen> createState() => _ActiveRideScreenState();
}

class _ActiveRideScreenState extends State<ActiveRideScreen> {
  final controller = Get.put(ActiveRideController());
  GoogleMapController? _mapController;
  final Set<Marker> _markers = {};
  final Set<Polyline> _polylines = {};

  @override
  void initState() {
    super.initState();
    controller.fetchActiveRide();
    // Listen to real-time updates via Socket.io
    controller.setupSocketListeners();
  }

  void _updateMapWithRide(Ride ride) {
    _markers.clear();
    _polylines.clear();

    final passenger = ride.passengers.first;

    // Pickup marker
    _markers.add(
      Marker(
        markerId: MarkerId('pickup'),
        position: LatLng(
          passenger.pickupLocation.lat,
          passenger.pickupLocation.lng,
        ),
        infoWindow: InfoWindow(title: 'Pickup'),
        icon: BitmapDescriptor.defaultMarkerWithHue(
          BitmapDescriptor.hueGreen,
        ),
      ),
    );

    // Dropoff marker
    _markers.add(
      Marker(
        markerId: MarkerId('dropoff'),
        position: LatLng(
          passenger.dropoffLocation.lat,
          passenger.dropoffLocation.lng,
        ),
        infoWindow: InfoWindow(title: 'Dropoff'),
        icon: BitmapDescriptor.defaultMarkerWithHue(
          BitmapDescriptor.hueRed,
        ),
      ),
    );

    // Route polyline
    _polylines.add(
      Polyline(
        polylineId: PolylineId('route'),
        points: [
          LatLng(
            passenger.pickupLocation.lat,
            passenger.pickupLocation.lng,
          ),
          LatLng(
            passenger.dropoffLocation.lat,
            passenger.dropoffLocation.lng,
          ),
        ],
        color: Colors.blue,
        width: 4,
      ),
    );

    setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Active Ride'),
      ),
      body: Obx(
        () {
          if (controller.isLoading.value) {
            return Center(child: CircularProgressIndicator());
          }

          final ride = controller.activeRide.value;

          if (ride == null) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text('No Active Ride'),
                  SizedBox(height: 16),
                  ElevatedButton(
                    onPressed: () => Get.offNamed('/driver/requests'),
                    child: Text('View Ride Requests'),
                  ),
                ],
              ),
            );
          }

          _updateMapWithRide(ride);
          final passenger = ride.passengers.first;

          return Column(
            children: [
              // Map
              Expanded(
                flex: 2,
                child: GoogleMap(
                  onMapCreated: (controller) => _mapController = controller,
                  initialCameraPosition: CameraPosition(
                    target: LatLng(
                      passenger.pickupLocation.lat,
                      passenger.pickupLocation.lng,
                    ),
                    zoom: 15,
                  ),
                  markers: _markers,
                  polylines: _polylines,
                ),
              ),

              // Ride details panel
              Expanded(
                flex: 1,
                child: SingleChildScrollView(
                  padding: EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      // Status badge
                      Container(
                        padding: EdgeInsets.symmetric(
                          horizontal: 12,
                          vertical: 8,
                        ),
                        decoration: BoxDecoration(
                          color: _getStatusColor(ride.status)
                              .withOpacity(0.2),
                          border: Border.all(
                            color: _getStatusColor(ride.status),
                          ),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Row(
                          children: [
                            Container(
                              width: 8,
                              height: 8,
                              decoration: BoxDecoration(
                                color: _getStatusColor(ride.status),
                                shape: BoxShape.circle,
                              ),
                            ),
                            SizedBox(width: 8),
                            Text(
                              _getStatusText(ride.status),
                              style: TextStyle(
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                          ],
                        ),
                      ),

                      SizedBox(height: 16),

                      // Passenger info
                      Text(
                        'Passenger',
                        style: TextStyle(fontWeight: FontWeight.bold),
                      ),
                      SizedBox(height: 8),
                      Card(
                        child: ListTile(
                          leading: CircleAvatar(child: Text('👤')),
                          title: Text(passenger.name),
                          subtitle: Text(_getPassengerStatus(
                            passenger.status,
                          )),
                        ),
                      ),

                      SizedBox(height: 16),

                      // Trip details
                      Text(
                        'Trip Details',
                        style: TextStyle(fontWeight: FontWeight.bold),
                      ),
                      SizedBox(height: 8),
                      _buildLocationItem(
                        'Pickup',
                        passenger.pickupLocation.address,
                        Colors.green,
                      ),
                      SizedBox(height: 8),
                      _buildLocationItem(
                        'Dropoff',
                        passenger.dropoffLocation.address,
                        Colors.red,
                      ),

                      SizedBox(height: 16),

                      // Fare
                      Container(
                        padding: EdgeInsets.all(12),
                        decoration: BoxDecoration(
                          color: Colors.grey.withOpacity(0.1),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Row(
                          mainAxisAlignment:
                              MainAxisAlignment.spaceBetween,
                          children: [
                            Text('Trip Fare'),
                            Text(
                              _formatCurrency(ride.totalFare),
                              style: TextStyle(
                                fontWeight: FontWeight.bold,
                                fontSize: 18,
                              ),
                            ),
                          ],
                        ),
                      ),

                      SizedBox(height: 16),

                      // Action buttons
                      ..._buildActionButtons(ride),
                    ],
                  ),
                ),
              ),
            ],
          );
        },
      ),
    );
  }

  Widget _buildLocationItem(
    String label,
    String address,
    Color color,
  ) {
    return Container(
      padding: EdgeInsets.all(12),
      decoration: BoxDecoration(
        border: Border.all(color: Colors.grey.withOpacity(0.3)),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        children: [
          Container(
            width: 10,
            height: 10,
            decoration: BoxDecoration(
              color: color,
              shape: BoxShape.circle,
            ),
          ),
          SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  label,
                  style: TextStyle(fontSize: 12, color: Colors.grey),
                ),
                Text(address),
              ],
            ),
          ),
        ],
      ),
    );
  }

  List<Widget> _buildActionButtons(Ride ride) {
    List<Widget> buttons = [];

    if (ride.status == 'accepted') {
      buttons.add(
        SizedBox(
          width: double.infinity,
          child: ElevatedButton(
            onPressed: () => controller.updateStatus('in-progress'),
            child: Text('Start Trip'),
          ),
        ),
      );
    } else if (ride.status == 'in-progress') {
      buttons.add(
        SizedBox(
          width: double.infinity,
          child: ElevatedButton(
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.green,
            ),
            onPressed: () => controller.updateStatus('completed'),
            child: Text('Complete Trip'),
          ),
        ),
      );
    }

    buttons.add(SizedBox(height: 8));
    buttons.add(
      SizedBox(
        width: double.infinity,
        child: ElevatedButton(
          style: ElevatedButton.styleFrom(
            backgroundColor: Colors.red,
          ),
          onPressed: () => controller.updateStatus('cancelled'),
          child: Text('Cancel Ride'),
        ),
      ),
    );

    return buttons;
  }

  Color _getStatusColor(String status) {
    switch (status) {
      case 'accepted':
        return Colors.blue;
      case 'in-progress':
        return Colors.orange;
      case 'completed':
        return Colors.green;
      case 'cancelled':
        return Colors.red;
      default:
        return Colors.grey;
    }
  }

  String _getStatusText(String status) {
    switch (status) {
      case 'accepted':
        return 'Going to Pickup';
      case 'in-progress':
        return 'Trip In Progress';
      default:
        return status.replaceFirstMapped(
          RegExp(r'^.'),
          (match) => match.group(0)!.toUpperCase(),
        );
    }
  }

  String _getPassengerStatus(String status) {
    switch (status) {
      case 'pending':
        return 'Waiting for pickup';
      case 'picked':
        return 'On board';
      case 'dropped':
        return 'Dropped off';
      default:
        return status;
    }
  }

  String _formatCurrency(double amount) {
    return '${amount.toStringAsFixed(2)} ETB';
  }

  @override
  void dispose() {
    _mapController?.dispose();
    super.dispose();
  }
}
```

---

## 📱 STATE MANAGEMENT: GetX Controllers

**File**: `lib/features/user/presentation/controllers/ride_booking_controller.dart`

```dart
import 'package:get/get.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:dio/dio.dart';

class RideBookingController extends GetxController {
  final Dio dio = Dio();
  
  var pickupLocation = Rx<LatLng?>(null);
  var dropoffLocation = Rx<LatLng?>(null);
  var wantPooling = false.obs;
  var isLoading = false.obs;
  var isJoiningPool = false.obs;
  var fareInfo = Rx(null);
  var nearbyDrivers = <Driver>[].obs;
  var poolInfo = Rx(null);

  @override
  void onInit() {
    super.onInit();
    // Check query params for pool join
    final poolId = Get.parameters['poolId'];
    if (poolId != null) {
      joinPool(poolId);
    }
  }

  void setPickupLocation(LatLng location) {
    pickupLocation.value = location;
    calculateFare();
  }

  void setDropoffLocation(LatLng location) {
    dropoffLocation.value = location;
    calculateFare();
  }

  void setPooling(bool value) {
    wantPooling.value = value;
    calculateFare();
  }

  Future<void> calculateFare() async {
    if (pickupLocation.value == null || dropoffLocation.value == null) {
      return;
    }

    try {
      final response = await dio.post(
        'http://localhost:3014/v1/pricing/calculate',
        data: {
          'ride_type': 'ECONOMY',
          'distance_meters': _calculateDistance(),
          'duration_seconds': 900, // 15 minutes estimate
          'is_pool': wantPooling.value,
          'active_rides': 50,
          'available_drivers': 20,
        },
      );

      fareInfo.value = response.data['data'];
      update();
    } catch (e) {
      print('Error calculating fare: $e');
    }
  }

  Future<void> findNearbyDrivers() async {
    if (pickupLocation.value == null) return;

    try {
      final response = await dio.get(
        'http://localhost:3011/v1/dispatch/nearby-drivers',
        queryParameters: {
          'latitude': pickupLocation.value!.latitude,
          'longitude': pickupLocation.value!.longitude,
          'radius': 5000, // 5km
        },
      );

      nearbyDrivers.value = (response.data['data'] as List)
          .map((d) => Driver.fromJson(d))
          .toList();
      update();
    } catch (e) {
      print('Error fetching nearby drivers: $e');
    }
  }

  Future<void> requestRide() async {
    if (pickupLocation.value == null || dropoffLocation.value == null) {
      throw Exception('Locations not set');
    }

    isLoading.value = true;
    try {
      final response = await dio.post(
        'http://localhost:3010/v1/rides/create',
        data: {
          'pickup_lat': pickupLocation.value!.latitude,
          'pickup_lng': pickupLocation.value!.longitude,
          'dropoff_lat': dropoffLocation.value!.latitude,
          'dropoff_lng': dropoffLocation.value!.longitude,
          'ride_type': 'ECONOMY',
          'want_pooling': wantPooling.value,
        },
      );

      // Navigate to tracking screen
      Get.offNamed(
        '/user/tracking',
        parameters: {'rideId': response.data['data']['ride_id']},
      );
    } catch (e) {
      throw Exception('Failed to request ride: $e');
    } finally {
      isLoading.value = false;
    }
  }

  Future<void> joinPool(String poolId) async {
    isJoiningPool.value = true;
    try {
      final response = await dio.get(
        'http://localhost:3013/v1/pooling/pool/$poolId',
      );

      poolInfo.value = response.data['data'];
      
      // Set locations
      final pool = response.data['data'];
      pickupLocation.value = LatLng(
        pool['pickup_location']['lat'],
        pool['pickup_location']['lng'],
      );
      dropoffLocation.value = LatLng(
        pool['dropoff_location']['lat'],
        pool['dropoff_location']['lng'],
      );
      
      wantPooling.value = true;
      update();
    } catch (e) {
      print('Error joining pool: $e');
    }
  }

  Future<void> joinExistingPool() async {
    if (poolInfo.value == null) {
      throw Exception('Pool info not available');
    }

    isLoading.value = true;
    try {
      await dio.post(
        'http://localhost:3013/v1/pooling/pool/${poolInfo.value['id']}/join',
        data: {
          'pickup_lat': pickupLocation.value!.latitude,
          'pickup_lng': pickupLocation.value!.longitude,
          'dropoff_lat': dropoffLocation.value!.latitude,
          'dropoff_lng': dropoffLocation.value!.longitude,
        },
      );
    } catch (e) {
      throw Exception('Failed to join pool: $e');
    } finally {
      isLoading.value = false;
    }
  }

  int _calculateDistance() {
    // Simple distance calculation (in production use geodistance package)
    final lat1 = pickupLocation.value!.latitude;
    final lng1 = pickupLocation.value!.longitude;
    final lat2 = dropoffLocation.value!.latitude;
    final lng2 = dropoffLocation.value!.longitude;

    const double R = 6371000; // Earth's radius in meters
    final dLat = (lat2 - lat1) * 3.14159 / 180;
    final dLng = (lng2 - lng1) * 3.14159 / 180;

    final a = (dLat / 2).sin() * (dLat / 2).sin() +
        (dLng / 2).sin() * (dLng / 2).sin();
    final c = 2 * (a.sqrt() / (1 - a).sqrt()).atan();

    return (R * c).toInt();
  }
}

class Driver {
  final String id;
  final String name;
  final String vehicleType;
  final double rating;
  final double distance;
  final LatLng location;

  Driver({
    required this.id,
    required this.name,
    required this.vehicleType,
    required this.rating,
    required this.distance,
    required this.location,
  });

  factory Driver.fromJson(Map<String, dynamic> json) {
    return Driver(
      id: json['id'],
      name: json['name'],
      vehicleType: json['vehicle_type'],
      rating: json['rating'].toDouble(),
      distance: json['distance'].toDouble(),
      location: LatLng(
        json['location']['lat'],
        json['location']['lng'],
      ),
    );
  }
}
```

---

## ✅ KEY DIFFERENCES: React ↔️ Flutter

| Aspect | React | Flutter |
|--------|-------|---------|
| **Maps** | Leaflet | Google Maps Flutter |
| **State** | Zustand | GetX |
| **HTTP** | Axios | Dio |
| **Real-time** | Socket.io | Socket.io-client |
| **Routing** | React Router | GetX Navigation |
| **Forms** | React Hook Form | TextFormField |
| **Styling** | TailwindCSS | Material Widgets |

---

## 📋 MIGRATION CHECKLIST

```
☐ Create Flutter project structure
☐ Create models (Ride, Driver, Location)
☐ Create GetX controllers
☐ Create presentation screens
☐ Implement Google Maps integration
☐ Implement Socket.io integration
☐ Connect to Go backend APIs
☐ Run end-to-end tests
☐ Optimize performance
☐ Deploy to TestFlight/Play Store
```

This guide provides complete code templates for migrating React components to Flutter widgets while maintaining functionality and adding mobile-specific enhancements.

