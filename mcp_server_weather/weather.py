def get_weather_data(city: str) -> dict:
    """Retrieve current weather conditions for a city."""
    return {
        "city": city,
        "temperature": 25,
        "condition": "Sunny",
        "humidity": 60,
        "wind_speed": 10
    }
