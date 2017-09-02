"""
Module operate Git repository.
"""
from git import Repo

class GitRepo:
    """
    Class operating Git repository.
    """    
    
    def __init__(self):
        """
        Constructor.
        """

    def open(self, path):
        """
        Open Git Repository.
        Arguments:
          path :: path to repository
        Return Values:
          True :: success
          False :: failure
        """
        try:
            repo = Repo(path)
        except:
            return False
        else:
            return True
        
    
    
