import cv2 as ocv
import colormath.color_constants as clc
import colormath.color_objects as clo
from colormath import color_diff_matrix
from colormath.color_conversions import convert_color
from colormath.color_diff import delta_e_cie2000, _get_lab_color1_vector, _get_lab_color2_matrix
from numba import njit, prange, jit
from numpy import ndarray, ubyte
import numpy as np
from sklearn.cluster import KMeans, MiniBatchKMeans
import matplotlib.pyplot as plt


color_palette = {
        "red": [convert_color(clo.sRGBColor(1, 0, 0), clo.LabColor), clo.sRGBColor(1, 0, 0)],
        "blue": [convert_color(clo.sRGBColor(0, 0, 1), clo.LabColor), clo.sRGBColor(0, 0, 1)],
        "green": [convert_color(clo.sRGBColor(0, 1, 0), clo.LabColor), clo.sRGBColor(0, 1, 0)],
        "pink": [convert_color(clo.sRGBColor(254,105,180, is_upscaled=True), clo.LabColor),
                 clo.sRGBColor(254,105,180, is_upscaled=True)],
        "brown": [convert_color(clo.sRGBColor(139,69,19, is_upscaled=True), clo.LabColor),
                  clo.sRGBColor(139,69,19, is_upscaled=True)],
        "tan": [convert_color(clo.sRGBColor(210,180,140, is_upscaled=True), clo.LabColor),
                clo.sRGBColor(210,180,140, is_upscaled=True)],
        "black": [convert_color(clo.sRGBColor(0, 0, 0), clo.LabColor), clo.sRGBColor(0, 0, 0)],
        "white": [convert_color(clo.sRGBColor(1, 1, 1), clo.LabColor), clo.sRGBColor(1, 1, 1)]
    }


def monkey_delta_e_cie2000(color1, color2, Kl=1, Kc=1, Kh=1):
    """
    Calculates the Delta E (CIE2000) of two colors.
    """
    color1_vector = _get_lab_color1_vector(color1)
    color2_matrix = _get_lab_color2_matrix(color2)
    delta_e: ndarray = color_diff_matrix.delta_e_cie2000(
        color1_vector, color2_matrix, Kl=Kl, Kc=Kc, Kh=Kh)[0]
    return delta_e.item()


# @jit(forceobj=True)
def pixel_color_snap_bgr(pix: list, palette: list[list]) -> list:
    pix_in_rgb: clo.LabColor = convert_color(clo.sRGBColor(pix[2], pix[1], pix[0], is_upscaled=True), clo.LabColor)
    # pix_in: clo.LabColor = convert_color(clo.sRGBColor(*pix, is_upscaled=True), clo.LabColor)

    color_out = color_palette["black"][1]
    # dist_from_black = monkey_delta_e_cie2000(pix_in, color_palette["black"][0])
    min_E: float = float("inf")
    # for v in color_palette.values():

    converted_palette_rgb = [(convert_color(clo.sRGBColor(c[2], c[1], c[0], is_upscaled=True), clo.LabColor), clo.sRGBColor(c[2], c[1], c[0], is_upscaled=True)) for c in palette]
    # converted_palette_rgb.append((convert_color(clo.sRGBColor(146, 48, 42, is_upscaled=True), clo.LabColor),
    #                               clo.sRGBColor(146, 48, 42, is_upscaled=True)))
    # converted_palette_rgb.append((convert_color(clo.sRGBColor(148, 111, 58, is_upscaled=True), clo.LabColor),
    #                               clo.sRGBColor(148, 111, 58, is_upscaled=True)))
    for c in converted_palette_rgb:
        delta_e = monkey_delta_e_cie2000(pix_in_rgb, c[0])
        if delta_e < min_E:
            min_E = delta_e
            color_out = c[1]

    # opencv uses BGR, for reasons
    return [color_out.rgb_b, color_out.rgb_g, color_out.rgb_r]
    # return color_out.get_value_tuple()


def kmeans_colors_rgb(img: ndarray):
    # https://medium.com/analytics-vidhya/color-separation-in-an-image-using-kmeans-clustering-using-python-f994fa398454
    # kmeans = KMeans(n_clusters=10)
    kmeans = MiniBatchKMeans(n_clusters=5)
    img = img.reshape((img.shape[1] * img.shape[0], 3))
    s = kmeans.fit(img)

    labels = kmeans.labels_
    print(labels)
    labels = list(labels)

    centroid = kmeans.cluster_centers_
    print(centroid)

    percent = []
    for i in range(len(centroid)):
        j = labels.count(i)
        j = j / (len(labels))
        percent.append(j)
    print(percent)

    plt.pie(percent, colors=np.array(centroid / 255), labels=np.arange(len(centroid)))
    plt.show()

    return centroid


# @jit(parallel=True, forceobj=True)
def array_loop_bgr(inpt: ndarray, palette: list) -> ndarray:
    cnvrt: ndarray = ndarray((inpt.shape[0], 3))
    # for i in range(resized_down.shape[0]):
    #     for j in range(resized_down.shape[1]):
    for i in prange(inpt.shape[0]):
        pix = pixel_color_snap_bgr(inpt[i], palette)
        cnvrt[i] = pix

    return cnvrt


def resize_test():
    image_bgr: ndarray = ocv.imread("red_zaku.jpg")
    original_width: int = image_bgr.shape[1]
    original_height: int = image_bgr.shape[0]

    down_width: int = original_width // 10
    down_height: int = original_height // 10
    down_points = (down_width, down_height)
    resized_down_bgr: ndarray = ocv.resize(image_bgr, down_points, interpolation=ocv.INTER_AREA)

    image_rgb: ndarray = ocv.cvtColor(image_bgr, ocv.COLOR_BGR2RGB)
    color_centroids_rgb = kmeans_colors_rgb(image_rgb)

    color_centroids_bgr = [[c[2], c[1], c[0]] for c in color_centroids_rgb]

    resized_down_bgr_1d: ndarray = resized_down_bgr.reshape((resized_down_bgr.shape[0] * resized_down_bgr.shape[1], 3))
    print("array loop")
    cc_bgr = array_loop_bgr(resized_down_bgr_1d, color_centroids_bgr)
    print("array loop done")
    cc_rgb = [[pix[2], pix[1], pix[0]] for pix in cc_bgr]
    color_resize = np.array(cc_rgb).reshape(resized_down_bgr.shape)

    up_points = (original_width, original_height)
    resized_up: ndarray = ocv.resize(color_resize, up_points, interpolation=ocv.INTER_AREA)

    plt.imshow(resized_up)
    plt.show()


if __name__ == "__main__":
    resize_test()
