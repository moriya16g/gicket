"""
usage: python -m unittest test.test_gitrepo
"""
import unittest
import gitrepo

class TestGitrepo(unittest.TestCase):
    
    def setUp(self):
        self.repo = gitrepo.GitRepo()

    def tearDown(self):
        pass

    def test_open(self):
        ret = self.repo.open("./")
        self.assertFalse(ret)

        ret = self.repo.open("../")
        self.assertTrue(ret)

