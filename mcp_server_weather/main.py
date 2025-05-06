from server import app
from weather import get_weather_data, weather_data_2

if __name__ == "__main__":
    app.add_tool(get_weather_data)
    app.add_tool(weather_data_2)
    app.run(transport="sse")
