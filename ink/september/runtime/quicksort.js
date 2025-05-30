export const sortBy = (v, pred) => (() => {
  let vPred = map(v, pred);
  let partition = (v, lo, hi) => (() => {
    let __ink_trampolined_lsub;
    let __ink_trampolined_rsub;
    let pivot = (() => {let __ink_acc_trgt = __as_ink_string(vPred); return __is_ink_string(__ink_acc_trgt) ? __ink_acc_trgt.valueOf()[(() => { return lo })()] || null : (__ink_acc_trgt[(() => { return lo })()] !== undefined ? __ink_acc_trgt[(() => { return lo })()] : null)})();
    let lsub = i => (() => { __ink_trampolined_lsub = i => __ink_match(((() => {let __ink_acc_trgt = __as_ink_string(vPred); return __is_ink_string(__ink_acc_trgt) ? __ink_acc_trgt.valueOf()[(() => { return i })()] || null : (__ink_acc_trgt[(() => { return i })()] !== undefined ? __ink_acc_trgt[(() => { return i })()] : null)})() < pivot), [[() => (true), () => (__ink_trampoline(__ink_trampolined_lsub, __as_ink_string(i + 1)))], [() => (false), () => (i)]]); return __ink_resolve_trampoline(__ink_trampolined_lsub, i) })();
    let rsub = j => (() => { __ink_trampolined_rsub = j => __ink_match(((() => {let __ink_acc_trgt = __as_ink_string(vPred); return __is_ink_string(__ink_acc_trgt) ? __ink_acc_trgt.valueOf()[(() => { return j })()] || null : (__ink_acc_trgt[(() => { return j })()] !== undefined ? __ink_acc_trgt[(() => { return j })()] : null)})() > pivot), [[() => (true), () => (__ink_trampoline(__ink_trampolined_rsub, (j - 1)))], [() => (false), () => (j)]]); return __ink_resolve_trampoline(__ink_trampolined_rsub, j) })();
    return (() => {
      let __ink_trampolined_sub;
      let sub;
      return sub = (i, j) => (() => {
        __ink_trampolined_sub = (i, j) => (() => {
          i = lsub(i);
          j = rsub(j);
          return __ink_match(
            (i < j), [
              [() => (false), () => (j)],
              [() => (true), () => ((() => {
                let tmp;
                let tmpPred;
                tmp = (() => {let __ink_acc_trgt = __as_ink_string(v); return __is_ink_string(__ink_acc_trgt) ? __ink_acc_trgt.valueOf()[(() => { return i })()] || null : (__ink_acc_trgt[(() => { return i })()] !== undefined ? __ink_acc_trgt[(() => { return i })()] : null)})();
                tmpPred = (() => {let __ink_acc_trgt = __as_ink_string(vPred); return __is_ink_string(__ink_acc_trgt) ? __ink_acc_trgt.valueOf()[(() => { return i })()] || null : (__ink_acc_trgt[(() => { return i })()] !== undefined ? __ink_acc_trgt[(() => { return i })()] : null)})();
                (() => {let __ink_assgn_trgt = __as_ink_string(v); __is_ink_string(__ink_assgn_trgt) ? __ink_assgn_trgt.assign((() => { return i })(), (() => {let __ink_acc_trgt = __as_ink_string(v); return __is_ink_string(__ink_acc_trgt) ? __ink_acc_trgt.valueOf()[(() => { return j })()] || null : (__ink_acc_trgt[(() => { return j })()] !== undefined ? __ink_acc_trgt[(() => { return j })()] : null)})()) : (__ink_assgn_trgt[(() => { return i })()]) = (() => {let __ink_acc_trgt = __as_ink_string(v); return __is_ink_string(__ink_acc_trgt) ? __ink_acc_trgt.valueOf()[(() => { return j })()] || null : (__ink_acc_trgt[(() => { return j })()] !== undefined ? __ink_acc_trgt[(() => { return j })()] : null)})(); return __ink_assgn_trgt})();
                (() => {let __ink_assgn_trgt = __as_ink_string(v); __is_ink_string(__ink_assgn_trgt) ? __ink_assgn_trgt.assign((() => { return j })(), tmp) : (__ink_assgn_trgt[(() => { return j })()]) = tmp; return __ink_assgn_trgt})();
                (() => {let __ink_assgn_trgt = __as_ink_string(vPred); __is_ink_string(__ink_assgn_trgt) ? __ink_assgn_trgt.assign((() => { return i })(), (() => {let __ink_acc_trgt = __as_ink_string(vPred); return __is_ink_string(__ink_acc_trgt) ? __ink_acc_trgt.valueOf()[(() => { return j })()] || null : (__ink_acc_trgt[(() => { return j })()] !== undefined ? __ink_acc_trgt[(() => { return j })()] : null)})()) : (__ink_assgn_trgt[(() => { return i })()]) = (() => {let __ink_acc_trgt = __as_ink_string(vPred); return __is_ink_string(__ink_acc_trgt) ? __ink_acc_trgt.valueOf()[(() => { return j })()] || null : (__ink_acc_trgt[(() => { return j })()] !== undefined ? __ink_acc_trgt[(() => { return j })()] : null)})(); return __ink_assgn_trgt})();
                (() => {let __ink_assgn_trgt = __as_ink_string(vPred); __is_ink_string(__ink_assgn_trgt) ? __ink_assgn_trgt.assign((() => { return j })(), tmpPred) : (__ink_assgn_trgt[(() => { return j })()]) = tmpPred; return __ink_assgn_trgt})();
                return __ink_trampoline(__ink_trampolined_sub, __as_ink_string(i + 1), (j - 1))
              })())]
            ])
        })();
        return __ink_resolve_trampoline(__ink_trampolined_sub, i, j)
      })()
    })()(lo, hi) })();
  return (() => {
    let __ink_trampolined_quicksort;
    let quicksort;
    return quicksort = (v, lo, hi) => (() => {
      __ink_trampolined_quicksort = (v, lo, hi) => __ink_match(len(v), [
        [() => (0), () => (v)],
        [() => (__Ink_Empty), () => (__ink_match((lo < hi), [
          [() => (false), () => (v)],
          [() => (true), () => ((() => {
            let p = partition(v, lo, hi);
            quicksort(v, lo, p);
            return __ink_trampoline(__ink_trampolined_quicksort, v, __as_ink_string(p + 1), hi)
          })())]
        ]))]
      ]);
      return __ink_resolve_trampoline(__ink_trampolined_quicksort, v, lo, hi)
    })()
  })()(v, 0, (len(v) - 1))
})();
export const sort__ink_em__ = v => sortBy(v, x => x);
export const sort = v => sort__ink_em__(clone(v));
