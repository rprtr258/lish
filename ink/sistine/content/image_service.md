#someday_maybe #project

https://github.com/rprtr258/Image-service
https://github.com/LaurentMazare/tch-rs/tree/main/examples/neural-style-transfer

remove dominant/red/green/blue color
http://www.aosabook.org/en/500L/making-your-own-image-filters.html#fnref1

make shaders parametrized:
    choose color
    several images
    floats

https://invata.handmade.network/
https://nothings.org/gamedev/blend_equations_with_destination_alpha.txt

migrate image to palette colors
[text to image](https://github.com/CompVis/stable-diffusion)
  https://github.com/h2non/imaginary
  https://github.com/cshum/imagor
  https://habr.com/ru/post/693512/

```embed
title: 'Стабильная диффузия для самых маленьких. Или строим свою собственную ярмарку с хороводом и скоморохами'
image: 'https://habrastorage.org/getpro/habr/upload_files/2aa/504/f8d/2aa504f8df431e97fa4dcde938bebdfd.jpg'
description: 'Волшебная сказка с лубочными картинками. Далеко ли, близко ли, высоко ли, низко ли, а летает нонче в небе жар-птица павлин из самого из города Муничинска. И где перо она потеряет, там картинки...'
url: 'https://habr.com/ru/post/709892/'
```

`stbimage.c`
```c
// https://nothings.org/gamedev/ssao/
#define STB_IMAGE_WRITE_IMPLEMENTATION
#include "stb_image_write.h"
#define STBI_NO_WRITE  // disable writer in old version of stb_image
#include "stb_image.c"
#define STB_PLOT_IMPLEMENTATION
#include "stb_plot.h" // unfinished, unreleased line graph library
#define STB_DEFINE
#include "stb.h"

int main(int argc, char **argv) {
   int i, n;
   char **data = stb_stringfile("c:/imv_log.txt", &n);
   for (i = 0; i < n; i += 2) {
      int w, h;
      int x0, y0, x1, y1, j, len;
      uint8 *pixels;
      char file1[999], file2[999], name[999];

      stbplot_dataset *ds = stbplot_dataset_create();
      stbplot_variable *v = stbplot_dependent_variable(ds, "brightness");

      if (sscanf(data[i], "%d%d%s", &x0, &y0, file1) != 3 || sscanf(data[i+1], "%d%d%s", &x1, &y1, file2) != 3) {
         stb_fatal("Error on line %d\n", i+1);
      }
      if (strcmp(file1, file2)) {
          stb_fatal("Mismatch %1 vs %2\n", i+1, i+2);
      }

      pixels = stbi_load(file1, &w, &h, NULL, 3);
      if (x0 >= w || y0 >= h || x1 >= w || y1 >= h) {
          stb_fatal("Bad point in %d/%d", i+1, i+2);
      }

      len = max(abs(x1 - x0), abs(y1 - y0));
      for (j = 0; j <= len; ++j) {
         int dx, dy, c, sum = 0;
         int x = (int) stb_linear_remap(j, 0, len, x0, x1);
         int y = (int) stb_linear_remap(j, 0, len, y0, y1);
         for (dx = -2; dx <= 2; ++dx) {
            for (dy = -2; dy <= 2; ++dy) {
               for (c = 0; c < 3; ++c) {
                  sum += pixels[((y + dy) * w + (x + dx)) * 3 + c];
               }
            }
         }
         stbplot_add_value(j, v, (double) sum / (9 * 3));
      }

      stb_splitpath(name, file1, STB_FILE);
      stbplot_plot(stb_sprintf("c:/temp2/graph_%s_%d.bmp", name, i / 2 + 1), ds, "Pixel brightness", 600, 400, 1);
      stbplot_dataset_destroy(ds);

      // overdraw the highlight
      for (j = 0; j <= len; ++j) {
         int dx, dy, c, color[3] = { 128,255,128 };
         int x = (int) stb_linear_remap(j,0,len,x0,x1);
         int y = (int) stb_linear_remap(j,0,len,y0,y1);
         for (dx = -2; dx <= 2; ++dx) {
            for (dy = -2; dy <= 2; ++dy) {
               for (c = 0; c < 3; ++c) {
                  pixels[((y + dy) * w + (x + dx)) * 3 + c] = color[c];
               }
            }
         }
      }
      stbi_write_png(stb_sprintf("c:/temp2/highlight_%s_%d.png", name, i / 2 + 1), w, h, 3, pixels, w * 3);
   }
   return 0;
}
```