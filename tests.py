#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from re import S
import shutil
import unittest
from pathlib import Path, PurePath

_test_dir = "_tests"

class MockFiles(unittest.TestCase):
    def __init__(self, name, file_list):
        self.name = name
        self.base_dir = PurePath(_test_dir).joinpath(name)
        for file in file_list:
            pp = self.base_dir.joinpath(file)
            Path(pp.parent).mkdir(parent=True, exist_ok=True)
            Path(pp).write_text(pp.name)
    
    def check_rename(self, filepath, new_name):
        pp = self.base_dir.joinpath(filepath)
        old_name = pp.name
        self.assertEqual(old_name, Path(pp.parent.joinpath(new_name)).read_text())

    def check_recovery(self, filepath):
        pp = self.base_dir.joinpath(filepath)
        self.assertEqual(pp.name, Path(pp).read_text())
    
    def tearDown(self):
        shutil.rmtree(str(self.base_dir))


class TestDownloadEdge(unittest.TestCase):

    def check_rename(self, filepath):
        pass
    
    def test_anime_rename(self):
        pass

    def test_recovery(self):
        pass

    def setUp(self):
        pass

    def teardown(self):
        pass


if __name__ == '__main__':
    unittest.main()