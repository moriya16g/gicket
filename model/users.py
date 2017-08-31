"""
Module handling user information.
"""

class User:
    """
    Class representing a user.
    """

    def __init__(self, name, mail_address):
        """
        Constructor.
        Arguments:
          name :: name used by Git
          mail_address :: mail address used by Git
        """
        self.name = name
        self.mail = mail_address

        
class Users:
    """
    Class representing collection of users.
    """

    def __init__(self):
        """
        Constructor.
        """
        self.users = []

    def add_user(self, user):
        """
        Add a user instans.
        Do not register users with duplicate mail addresses.
        Arguments:
          user :: instans of User class
        Return Values:
          True :: add user
          False :: do not add
        """
        for i in range(0, len(self.users)):
            if self.users[i].mail == user.mail:
                return False
        self.users.append(user)
        return True

    def get_user(self, index):
        if len(self.users) <= index:
            return "",""
        return self.users[index].name, self.users[index].mail
    
