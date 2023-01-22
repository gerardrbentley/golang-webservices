import streamlit as st
import httpx
import os
import json

BACKEND_HOST = os.getenv("BACKEND_HOST", "http://127.0.0.1:5000")

st.header("Search Seattle Emergency Food Locations")

query = st.text_input("Search term")
search_url = f"{BACKEND_HOST}/place"

response = httpx.get(search_url, params={"name": query})
try:
    records = response.json()
    st.dataframe(records)
except (json.decoder.JSONDecodeError, st.errors.StreamlitAPIException):
    st.warning(f"Error from Web Request Code: {response.status_code}")
    st.stop()

map_data = []
for record in records:
    try:
        record["latitude"] = float(record["latitude"])
        record["longitude"] = float(record["longitude"])
        map_data.append(record)
    except ValueError:
        pass

st.map(map_data)
