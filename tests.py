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

class MockFiles(unittest.TestCase):
    def __init__(self, name, file_list):
        super().__init__()
        self.name = name
        self.base_dir = PurePath(_test_dir).joinpath(name)
        for file in file_list:
            pp = self.base_dir.joinpath(file)
            Path(pp.parent).mkdir(parents=True, exist_ok=True)
            Path(pp).write_text(pp.name)
    
    def check_rename(self, filepath, new_name, new_dir=None):
        pp = self.base_dir.joinpath(filepath)
        old_name = pp.name
        if new_dir is None:
            new_path = pp.parent
        else:
            new_path = self.base_dir.joinpath(new_dir)
        self.assertEqual(old_name, Path(new_path.joinpath(new_name)).read_text())

    def check_recovery(self, filepath):
        pp = self.base_dir.joinpath(filepath)
        self.assertEqual(pp.name, Path(pp).read_text())

    def check_exists(self, filepath):
        self.assertTrue(Path(self.base_dir.join(filepath)).exists())
    
_mock_files_tvdb = [
    "1-201/银魂.Gintama.003.mp4",
    "1-201/银魂.Gintama.003.chs.ass",
    "S2/银魂.Gintama.202.mkv", 
    "S2/银魂.Gintama.202.ass", 
    "S2/银魂.Gintama.203.ass",
    "op/海贼王第10集.mkv",
]

class TestAnimeRenamer(unittest.TestCase):

    def run_app(self, options, base_dir, config=None, custom=None):
        
        def write_config(cfg, name):
            p = base_dir.joinpath(name)
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

        try:
            return subprocess.check_output(cmds, stderr=subprocess.STDOUT)
        except subprocess.CalledProcessError as e:
            self.fail(e.output)
    
    def test_anime_rename_tvdb(self):
        renamed_files = [
            "Gintama.S01E03.[003].mp4",
            "Gintama.S01E03.[003].chs.ass",
            "Gintama.S05E01.[202].mkv",
            "Gintama.S05E01.[202].ass",
            "银魂.Gintama.203.ass",
            "One Piece.S02E02.[010].mkv"
        ]
        mock_fs = MockFiles("anime_ranme", _mock_files_tvdb)
        options = [
            "--language", "en",
            "--db", "tvdb",
            "-p", r"银魂\.(?P<series>.+)\.(?P<absolute>\d+)\.(?P<ext>\w+?)$",
            "-p", r"(?P<series>海贼王)第(?P<absolute>\d+)集\.(?P<ext>\w+?)$",
            "--format", r"{series}.S{season.2}E{episode.2}.[{absolute.3}].{ext}",
            str(mock_fs.base_dir)
        ]
        o = self.run_app(options, mock_fs.base_dir)
        # print(o)
        for i in range(0, len(_mock_files_tvdb)):
            mock_fs.check_rename(_mock_files_tvdb[i], renamed_files[i])

    def test_regexp_rename(self):
        renamed_files = [
            "Gintama.[003].mp4",
            "Gintama.[003].chs.ass",
            "Gintama.[202].mkv",
            "Gintama.[202].ass",
            "银魂.Gintama.203.ass",
            "海贼王.[010].mkv"
        ]
        mock_fs = MockFiles("regexp_ranme", _mock_files_tvdb)
        options = [
            "-R",
            "-p", r"银魂\.(?P<name>.+)\.(?P<absolute>\d+)\.\w+?$",
            "-p", r"(?P<name>海贼王)第(?P<absolute>\d+)集\.\w+?$",
            "--format", r"{name}.[{absolute.3}].{ext}",
            str(mock_fs.base_dir)
        ]
        o = self.run_app(options, mock_fs.base_dir)
        for i in range(0, len(_mock_files_tvdb)):
            mock_fs.check_rename(_mock_files_tvdb[i], renamed_files[i])

    def test_recovery(self):
        mock_fs = MockFiles("recovery_ranme", _mock_files_tvdb)
        rename_options = [
            "-R",
            "-p", r"银魂\.(?P<name>.+)\.(?P<absolute>\d+)\.\w+?$",
            "-p", r"(?P<name>海贼王)第(?P<absolute>\d+)集\.\w+?$",
            "--format", r"{name}.[{absolute.3}].{ext}",
            str(mock_fs.base_dir)
        ]
        o = self.run_app(rename_options, mock_fs.base_dir)
        recovery_options = [
            "recovery", mock_fs.base_dir.joinpath("rename.1.log")
        ]
        o = self.run_app(recovery_options, mock_fs.base_dir)
        for f in _mock_files_tvdb:
            mock_fs.check_recovery(f)

    def test_move_to_dir(self):
        renamed_files = [
            ("Gintama.S01E03.[003].mp4", "Gintama-S01"),
            ("Gintama.S01E03.[003].chs.ass", "Gintama-S01"),
            ("Gintama.S05E01.[202].mkv", "Gintama-S05"),
            ("Gintama.S05E01.[202].ass", "Gintama-S05"),
            ("银魂.Gintama.203.ass", None),
            ("One Piece.S02E02.[010].mkv", "One Piece-S02"),
        ]
        mock_fs = MockFiles("move_to", _mock_files_tvdb)
        options = [
            "--language", "en",
            "--db", "tvdb",
            "-p", r"银魂\.(?P<series>.+)\.(?P<absolute>\d+)\.(?P<ext>\w+?)$",
            "-p", r"(?P<series>海贼王)第(?P<absolute>\d+)集\.(?P<ext>\w+?)$",
            "--format", r"{series}.S{season.2}E{episode.2}.[{absolute.3}].{ext}",
            "-m", "{series}-S{season.2}",
            str(mock_fs.base_dir)
        ]
        o = self.run_app(options, mock_fs.base_dir)
        # print(o)
        for i in range(0, len(_mock_files_tvdb)):
            mock_fs.check_rename(_mock_files_tvdb[i], renamed_files[i][0], renamed_files[i][1])

    def test_tmdb(self):
        _files = [
            ("one.piece.10.mkv", "One Piece.S01E10.[010].mkv"),
            ("one.piece.97.mkv", "One Piece.S04E97.[097].mkv"),
            ("one.piece.1022.mkv", "One Piece.S21E1022.[1022].mkv"),
        ]
        mock_fs = MockFiles("tmdb", [x[0] for x in _files])
        options = [
            "--language", "en",
            "--id", "37854",
            "--tmdb.absoluteGroupSeason", "Season 1 (Absolute Order)",
            "-p", r".*\.(?P<absolute>\d+)\.(?P<ext>\w+?)$",
            "--format", r"{series}.S{season.2}E{episode.2}.[{absolute.3}].{ext}",
            str(mock_fs.base_dir)
        ]
        o = self.run_app(options, mock_fs.base_dir)
        # print(o)
        for f in _files:
            mock_fs.check_rename(f[0], f[1])
        
    def setup(self):
        print ("setUp")

    def tearDown(self):
        try:
            shutil.rmtree(_test_dir)
            pass
        except:
            pass

if __name__ == '__main__':
    unittest.main()