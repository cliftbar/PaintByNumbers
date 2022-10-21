import numpy as np
from matplotlib import pyplot as plt
import skimage


def display(img, colorbar=False):
    "Displays an image."
    plt.figure(figsize=(10, 10))
    if len(img.shape) == 2:
        i = skimage.io.imshow(img, cmap='gray')
    else:
        i = skimage.io.imshow(img)
    if colorbar:
        plt.colorbar(i, shrink=0.5, label='depth')
    plt.tight_layout()
    plt.show()


def make_pattern(shape=(16, 16), levels=64):
    "Creates a pattern from gray values."
    return np.random.randint(0, levels - 1, shape) / levels


def create_circular_depthmap(shape=(600, 800), center=None, radius=100):
    "Creates a circular depthmap, centered on the image."
    depthmap = np.zeros(shape, dtype=np.float)
    r = np.arange(depthmap.shape[0])
    c = np.arange(depthmap.shape[1])
    R, C = np.meshgrid(r, c, indexing='ij')
    if center is None:
        center = np.array([r.max() / 2, c.max() / 2])
    d = np.sqrt((R - center[0])**2 + (C - center[1])**2)
    depthmap += (d < radius)
    return depthmap


def normalize(depthmap):
    "Normalizes values of depthmap to [0, 1] range."
    if depthmap.max() > depthmap.min():
        return (depthmap - depthmap.min()) / (depthmap.max() - depthmap.min())
    else:
        return depthmap


def make_autostereogram(depthmap, pattern, shift_amplitude=0.1, invert=False):
    "Creates an autostereogram from depthmap and pattern."
    depthmap = normalize(depthmap)
    if invert:
        depthmap = 1 - depthmap
    autostereogram = np.zeros_like(depthmap, dtype=pattern.dtype)
    for row in np.arange(autostereogram.shape[0]):  # for each row
        for col in np.arange(autostereogram.shape[1]):
            if col < pattern.shape[1]:
                autostereogram[row, col] = pattern[row % pattern.shape[0], col]
            else:
                shift = int(depthmap[row, col] * shift_amplitude * pattern.shape[1])
                autostereogram[row, col] = autostereogram[row, col - pattern.shape[1] + shift]
    return autostereogram


def main():
    pattern = make_pattern(shape=(128, 64))
    depthmap = create_circular_depthmap(radius=150)
    autostereogram = make_autostereogram(depthmap, pattern)
    display(autostereogram)


if __name__ == "__main__":
    main()
