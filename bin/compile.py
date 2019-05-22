import os, sys
GOOS = ['windows', 'linux', 'darwin', 'freebsd']
GOARCH = ['amd64', '386']
packname = sys.argv[1:]
print('[Start Compile]')

if not os.path.exists('target'):
    os.mkdir('target')

os.chdir('target')
path = os.getcwd()

for o in GOOS:
    os.environ['GOOS'] = o
    ARCHS = GOARCH.copy()
    if o in ['linux']:
        ARCHS += ['mips', 'mips64', 'mipsle', 'mips64le', 'arm', 'arm64']
    if o in ['freebsd']:
        ARCHS += ['arm']
    for a in ARCHS:
        os.environ['GOARCH'] = a
        for p in packname:
            target = '{os}_{arch}'.format(os=o, arch=a)
            print('Compiling', p)
            print('Compile for', o, a)
            if not os.path.exists(target):
                os.mkdir(target)
            os.chdir(target)
            os.system('go build ' + p)
            os.chdir(path)
            print('Finish\n')
