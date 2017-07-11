#!/bin/python
# -*- coding: utf-8 -*-

import os
import shutil
from os import walk
import re
import argparse

class Config(object):
    def __init__(self):
        self.malcolm_namespace = "malcolm"
        self.malcolm_address = "malcolm.youringress.com"
        self.malcolm_image = "u2takey/malcolm:v0.1"
        self.malcolm_mongo_image = "mongo:3.4.6"
        self.mogno_pvc_size = "20Gi"
        self.mogno_storage_class = "ceph"
      
    
    def getval(self, key):
        return self.__dict__[key]


variable = re.compile(r'{{.+?}}')
detail = re.compile(r'((\d+) )?([a-zA-Z_0-9-]+)')
def render_template(tmpl, config):
    matches = variable.findall(tmpl)
    for match in matches:
        segs = detail.search(match)
        if segs.group() == '':
            raise Exception('Error: Invalid template item(' + match + ')')
        value = config.getval(segs.group(3))
        spaces = segs.group(2)
        if spaces != '' and spaces != None:
            leading = ''.join(' ' for i in range(int(spaces)))
            value = str(value).replace('\n', '\n' + leading)
        tmpl = tmpl.replace(match, value)
    return tmpl

def generate(src, dst, config):
    shutil.rmtree(dst)
    shutil.copytree(src, dst)
    for (dirpath, dirnames, filenames) in walk(dst):
        for filename in filenames:
            if not filename.startswith('.'):
                with open(os.path.join(dirpath, filename), "r+") as f:
                    content = f.read()
                    f.seek(0)
                    content_new = render_template(content, config)
                    f.write(content_new)
                    f.truncate()

if __name__ == '__main__':
    base_dir = os.path.dirname(os.path.abspath(__file__))
    src_dir = os.path.join(base_dir, 'template')
    output_dir = os.path.join(base_dir, 'generate')
    if not os.path.exists(output_dir):
        os.makedirs(output_dir)
    parser = argparse.ArgumentParser(description='deploy utils.')
    args = parser.parse_args()
    print("generating deploy yaml....")
    config = Config()        
    generate(src_dir, output_dir, config)
    print("generating deploy yaml done")
    
    
