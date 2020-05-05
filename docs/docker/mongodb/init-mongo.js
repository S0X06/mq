
db.createUser(
	{
		user  : 'admin',
		pwd   : 'admin0725',
		roles : [ 
			{ 
				role:'userAdminAnyDatabase',
				db: 'news'
			}
		]
	}
);