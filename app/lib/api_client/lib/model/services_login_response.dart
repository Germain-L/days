//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class ServicesLoginResponse {
  /// Returns a new [ServicesLoginResponse] instance.
  ServicesLoginResponse({
    this.token,
    this.user,
  });

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? token;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  ServicesUserResponse? user;

  @override
  bool operator ==(Object other) => identical(this, other) || other is ServicesLoginResponse &&
    other.token == token &&
    other.user == user;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (token == null ? 0 : token!.hashCode) +
    (user == null ? 0 : user!.hashCode);

  @override
  String toString() => 'ServicesLoginResponse[token=$token, user=$user]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (this.token != null) {
      json[r'token'] = this.token;
    } else {
      json[r'token'] = null;
    }
    if (this.user != null) {
      json[r'user'] = this.user;
    } else {
      json[r'user'] = null;
    }
    return json;
  }

  /// Returns a new [ServicesLoginResponse] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static ServicesLoginResponse? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "ServicesLoginResponse[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "ServicesLoginResponse[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return ServicesLoginResponse(
        token: mapValueOfType<String>(json, r'token'),
        user: ServicesUserResponse.fromJson(json[r'user']),
      );
    }
    return null;
  }

  static List<ServicesLoginResponse> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <ServicesLoginResponse>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = ServicesLoginResponse.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, ServicesLoginResponse> mapFromJson(dynamic json) {
    final map = <String, ServicesLoginResponse>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = ServicesLoginResponse.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of ServicesLoginResponse-objects as value to a dart map
  static Map<String, List<ServicesLoginResponse>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<ServicesLoginResponse>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = ServicesLoginResponse.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
  };
}

