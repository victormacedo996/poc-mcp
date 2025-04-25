from server import app
from weather import get_weather_data

if __name__ == "__main__":
    app.add_tool(get_weather_data)
    app.run(transport="sse")
