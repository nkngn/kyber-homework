# System design problem
## Step 1: Outline use cases, constraints, and assumptions
### Use cases (Functional Requirements)
- Order Book Fetching and Processing
    - Continuously retrieve order book data from a specified exchange.
    - Store and update order book data dynamically.
- Optimal Trade Route Calculation
    - **REST API Endpoint**: Expose an endpoint that allows users to query the best trade route and price. Input: starting token (`Token Y`), target token (`Token X`), trade amount (`n` units of Token X). Output: The lowest effective ask price when buying n units of Token X và The highest effective bid price when selling n units of Token X.
    - **Multi-Hop Trading Support**: Consider multiple market pairs if a direct trading pair does not exist.
### Constraints and assumptions (Non-Functional Requirements)
- **Low Latency**: The system should provide near real-time responses for trade calculations.
- **High Availability**: The system should be resilient to API failures from exchanges.
- **Scalability**: The architecture should handle multiple exchanges and large trading volumes.
- **Fault Tolerance**: Ensure fallback mechanisms in case of API failures or data inconsistencies.
- **Security**: Secure API endpoints and prevent abuse (e.g., rate limiting).