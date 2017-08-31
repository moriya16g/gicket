"""
usage: % python -m unittest test.test_users
"""

import unittest
import users

class TestUsers(unittest.TestCase):
    
    def setUp(self):
        self.users = users.Users()

    def tearDown(self):
        pass

    def test_add_user(self):
        user1 = users.User("User 1", "user1.mailaddress.jp")
        self.users.add_user(user1)
        name, mail = self.users.get_user(0)
        self.assertEqual(name, "User 1")
        self.assertEqual(mail, "user1.mailaddress.jp")

        user2 = users.User("User 2", "user2.mailaddress.jp")
        self.users.add_user(user2)
        name, mail = self.users.get_user(0)
        self.assertEqual(name, "User 1")
        self.assertEqual(mail, "user1.mailaddress.jp")
        name, mail = self.users.get_user(1)
        self.assertEqual(name, "User 2")
        self.assertEqual(mail, "user2.mailaddress.jp")

        user3 = users.User("User 3", "user2.mailaddress.jp")
        self.users.add_user(user3)
        name, mail = self.users.get_user(0)
        self.assertEqual(name, "User 1")
        self.assertEqual(mail, "user1.mailaddress.jp")
        name, mail = self.users.get_user(1)
        self.assertEqual(name, "User 2")
        self.assertEqual(mail, "user2.mailaddress.jp")
        name, mail = self.users.get_user(2)
        self.assertEqual(name, "")
        self.assertEqual(mail, "")

        user4 = users.User("User 4", "user4.mailaddress.jp")
        self.users.add_user(user4)
        name, mail = self.users.get_user(0)
        self.assertEqual(name, "User 1")
        self.assertEqual(mail, "user1.mailaddress.jp")
        name, mail = self.users.get_user(1)
        self.assertEqual(name, "User 2")
        self.assertEqual(mail, "user2.mailaddress.jp")
        name, mail = self.users.get_user(2)
        self.assertEqual(name, "User 4")
        self.assertEqual(mail, "user4.mailaddress.jp")
        name, mail = self.users.get_user(3)
        self.assertEqual(name, "")
        self.assertEqual(mail, "")

