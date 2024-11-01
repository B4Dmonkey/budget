/* 
	# Budget
	This package contains the business logic for the budget application.

	The main object in this package is the Budget object.
	Currently this is a convenience object for ensuring that methods have access to the context and database connection.
	
	The key commands are as follows:
	- `NewBudget` - creates a new Budget object
	- `AddNewTransactionsFromDocument` - adds new transactions to the database from a CSV file. 
		This is currently implemented for Chase debit transactions.
*/
package budget