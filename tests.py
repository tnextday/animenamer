#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import subprocess
from re import S
import shutil
import unittest
from pathlib import Path, PurePath
import json

_test_dir = "_tests"
_app_path = "build/animenamer"
_tvdb_apikey = "Z5SC1ZD07NNS8TDC"

class MockFiles(unittest.TestCase):
    def __init__(self, name, file_list):
        super().__init__()
        self.name = name
        self.base_dir = PurePath(_test_dir).joinpath(name)
        for file in file_list:
            pp = self.base_dir.joinpath(file)
            Path(pp.parent).mkdir(parents=True, exist_ok=True)
            Path(pp).write_text(pp.name)
    
    def check_rename(self, filepath, new_name):
        pp = self.base_dir.joinpath(filepath)
        old_name = pp.name
        self.assertEqual(old_name, Path(pp.parent.joinpath(new_name)).read_text())

    def check_recovery(self, filepath):
        pp = self.base_dir.joinpath(filepath)
        self.assertEqual(pp.name, Path(pp).read_text())

    def check_exists(self, filepath):
        self.assertTrue(Path(self.base_dir.join(filepath)).exists())
    
    def tearDown(self):
        shutil.rmtree(str(self.base_dir))
    

class TestAnimeRenamer(unittest.TestCase):

    def run_app(self, options, mock_fs, config=None, custom=None):
        
        def write_config(cfg, name):
            p = mock_fs.base_dir.joinpath(name)
            Path(p).write_text(json.dumps(cfg, indent=4))
            return str(p)
        
        cmds = [_app_path] 
        # cmds += ["-v"]

        if config is not None:
            cmds += ["-c", write_config(config, "animenamer.json")]       
        if custom is not None:
            cmds += ["--custom", write_config(custom, "animenamer.custom.json")]       

        cmds += []
        cmds += options
        cmds += [str(mock_fs.base_dir)]

        try:
            return subprocess.check_output(cmds, stderr=subprocess.STDOUT)
        except subprocess.CalledProcessError as e:
            self.fail(e.output)
    
    def test_anime_rename(self):
        files = [
            ("1-201.BDRIP.720P.X264-10bit_AAC/银魂.Gintama.003.mp4","Gintama.S01E03.[003].mp4"),
            ("1-201.BDRIP.720P.X264-10bit_AAC/银魂.Gintama.003.chs.ass","Gintama.S01E03.[003].chs.ass"),
            ("银魂第二季.1080p.x264_AAC/银魂.Gintama.202.mkv", "Gintama.S05E01.[202].mkv"),
            ("银魂第二季.1080p.x264_AAC/银魂.Gintama.202.ass", "Gintama.S05E01.[202].ass"),
            ("银魂第二季.1080p.x264_AAC/银魂.Gintama.203.ass", "银魂.Gintama.203.ass"),
            ("op/海贼王第10集.mkv", "One Piece.S02E02.[010].mkv")
        ]
        mock_fs = MockFiles("anime_ranme", [x[0] for x in files])
        options = [
            "--apikey", _tvdb_apikey,
            "--language", "en",
            "-p", r"银魂\.(?P<series>.+)\.(?P<absolute>\d+)\.(?P<ext>\w+?)$",
            "-p", r"(?P<series>海贼王)第(?P<absolute>\d+)集\.(?P<ext>\w+?)$",
            "--format", r"{series}.S{season.2}E{episode.2}.[{absolute.3}].{ext}",
        ]
        o = self.run_app(options, mock_fs)
        # print(o)
        for r in files:
            mock_fs.check_rename(r[0], r[1])

    def test_regexp_rename(self):
        files = [
            ("1-201.BDRIP.720P.X264-10bit_AAC/银魂.Gintama.003.mp4","Gintama.[003].mp4"),
            ("1-201.BDRIP.720P.X264-10bit_AAC/银魂.Gintama.003.chs.ass","Gintama.[003].chs.ass"),
            ("银魂第二季.1080p.x264_AAC/银魂.Gintama.202.mkv", "Gintama.[202].mkv"),
            ("银魂第二季.1080p.x264_AAC/银魂.Gintama.202.ass", "Gintama.[202].ass"),
            ("银魂第二季.1080p.x264_AAC/银魂.Gintama.203.ass", "银魂.Gintama.203.ass"),
            ("op/海贼王第10集.mkv", "海贼王.[010].mkv")
        ]
        mock_fs = MockFiles("regexp_ranme", [x[0] for x in files])
        options = [
            "-R",
            "-p", r"银魂\.(?P<name>.+)\.(?P<absolute>\d+)\.\w+?$",
            "-p", r"(?P<name>海贼王)第(?P<absolute>\d+)集\.\w+?$",
            "--format", r"{name}.[{absolute.3}].{ext}",
        ]
        o = self.run_app(options, mock_fs)
        for r in files:
            mock_fs.check_rename(r[0], r[1])

    def test_recovery(self):
        pass

    def setUp(self):
        pass

    def teardown(self):
        # shutil.rmtree(_test_dir)
        pass


if __name__ == '__main__':
    unittest.main()