"""
Module operate vcs repository.
"""
import gitrepo

class Repository:
    """
    Class operating repository.
    """

    GIT_REPOSITORY = "Git"
    SVN_REPOSITORY = "Subversion"

    def __init__(self):
        """
        Constructor.
        """
    
    def open(self, vcs, path):
        """
        Open repository.
        Arguments:
          vcs :: version controll system name(only "Git") #ToDo
          path :: path to repository
        """
        if vcs == GIT_REPOSITORY:
            self.vcs = GIT_REPOSITORY
        else:
            self.vcs = GIT_REPOSITORY # default
        self.path = path
        

    def open_git():
        """
        Open Git Repository.
        Arguments:
          none
        Return Values:
          True :: success
          False :: failure
        """
        
