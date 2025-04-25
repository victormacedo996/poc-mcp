import sqlite3
import os


local_dir = os.path.dirname(os.path.abspath(__file__))
connection = sqlite3.connect(f"{local_dir}/gateway.db")

def create_tables():
    cursor = connection.cursor()
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS mcp_servers (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        description TEXT NOT NULL,
        name TEXT NOT NULL,
        url TEXT NOT NULL,
        environment TEXT NOT NULL,               
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    ''')
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS llm_models (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        size TEXT,                            
        num_parameters INTEGER,               
        description TEXT,                     
        pros TEXT,                            
        cons TEXT,                            
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    ''')
    connection.commit()
    cursor.close()

def get_available_mcp_servers():
    cursor = connection.cursor()
    cursor.execute('''
        SELECT url, name, description FROM mcp_servers WHERE environment = 'production'
    ''')
    rows = cursor.fetchall()
    cursor.close()
    return rows

def add_mcp_server(url, environment, name, description):
    cursor = connection.cursor()
    cursor.execute('''
        INSERT INTO mcp_servers (url, environment, name, description)
        VALUES (?, ?, ?, ?)
    ''', (url, environment, name, description))
    connection.commit()
    cursor.close()
    

def get_available_llms():
    cursor = connection.cursor()
    cursor.execute('''
        SELECT name, size, num_parameters, description, pros, cons FROM llm_models;
    ''')
    rows = cursor.fetchall()
    cursor.close()
    return rows